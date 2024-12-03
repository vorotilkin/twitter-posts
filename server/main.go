package main

import (
	"context"
	"github.com/vorotilkin/twitter-posts/infrastructure/repositories/comment"
	"github.com/vorotilkin/twitter-posts/infrastructure/repositories/like"
	"github.com/vorotilkin/twitter-posts/infrastructure/repositories/post"
	"github.com/vorotilkin/twitter-posts/interfaces"
	"github.com/vorotilkin/twitter-posts/pkg/configuration"
	"github.com/vorotilkin/twitter-posts/pkg/database"
	pkgGrpc "github.com/vorotilkin/twitter-posts/pkg/grpc"
	"github.com/vorotilkin/twitter-posts/pkg/migration"
	"github.com/vorotilkin/twitter-posts/proto"
	"github.com/vorotilkin/twitter-posts/usecases"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type config struct {
	Grpc struct {
		Server pkgGrpc.Config
	}
	Db        database.Config
	Migration migration.Config
}

func newConfig(configuration *configuration.Configuration) (*config, error) {
	c := new(config)
	err := configuration.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func main() {
	opts := []fx.Option{
		fx.Provide(zap.NewProduction),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(configuration.New),
		fx.Provide(newConfig),
		fx.Provide(func(c *config) pkgGrpc.Config {
			return c.Grpc.Server
		}),
		fx.Provide(func(c *config) database.Config {
			return c.Db
		}),
		fx.Provide(database.New),
		fx.Provide(func(c *config) migration.Config { return c.Migration }),
		fx.Provide(fx.Annotate(func(c *config) string { return c.Db.PostgresDSN() }, fx.ResultTags(`name:"dsn"`))),
		fx.Provide(fx.Annotate(pkgGrpc.NewServer,
			fx.As(new(grpc.ServiceRegistrar)),
			fx.As(new(interfaces.Hooker)))),
		fx.Provide(fx.Annotate(post.NewRepository, fx.As(new(usecases.PostsRepository)))),
		fx.Provide(fx.Annotate(like.NewRepository, fx.As(new(usecases.LikeRepository)))),
		fx.Provide(fx.Annotate(comment.NewRepository, fx.As(new(usecases.CommentRepository)))),
		fx.Provide(fx.Annotate(usecases.NewPostsServer, fx.As(new(proto.PostsServer)))),
		fx.Invoke(func(lc fx.Lifecycle, server interfaces.Hooker) {
			lc.Append(fx.Hook{
				OnStart: server.OnStart,
				OnStop:  server.OnStop,
			})
		}),
		fx.Invoke(fx.Annotate(migration.Do, fx.ParamTags("", "", `name:"dsn"`))),
		fx.Invoke(proto.RegisterPostsServer),
	}

	app := fx.New(opts...)
	err := app.Start(context.Background())
	if err != nil {
		panic(err)
	}

	<-app.Done()

	err = app.Stop(context.Background())
	if err != nil {
		panic(err)
	}
}
