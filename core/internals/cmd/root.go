package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/deformal/mongo-auto-type-gen/internals/config"
	"github.com/deformal/mongo-auto-type-gen/internals/infer"
	"github.com/deformal/mongo-auto-type-gen/internals/mongo"
	"github.com/deformal/mongo-auto-type-gen/internals/render"
	"github.com/deformal/mongo-auto-type-gen/pkg"
	"github.com/spf13/cobra"
)

type Options struct {
	URI               string
	Out               string
	Sample            int
	OptionalThreshold float64
	DateAs            string
	ObjectIDAs        string
	ConfigPath        string
	EnvFile           string
}

var opts Options

var rootCmd = &cobra.Command{
	Use:   "mongots",
	Short: "Generate TypeScript types from MongoDB collections by inference",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Config")
		fmt.Println("Sample", opts.Sample)
		if err := config.LoadDotEnv(opts.EnvFile); err != nil {
			return err
		}

		env := config.ReadEnv()

		if opts.URI == "" {
			opts.URI = env.MongoURI
		}
		if opts.Out == "" {
			opts.Out = env.Out
		}
		if opts.Sample == 200 && env.Sample != 200 {
			opts.Sample = env.Sample
		}
		if opts.OptionalThreshold == 0.98 && env.OptionalThreshold != 0.98 {
			opts.OptionalThreshold = env.OptionalThreshold
		}
		if opts.DateAs == "string" && env.DateAs != "" {
			opts.DateAs = env.DateAs
		}
		if opts.ObjectIDAs == "string" && env.ObjectIDAs != "" {
			opts.ObjectIDAs = env.ObjectIDAs
		}

		ctx := context.Background()

		client, err := mongo.Connect(ctx, opts.URI)
		if err != nil {
			return err
		}
		defer client.Disconnect(ctx)

		dbs, err := mongo.ListDatabases(ctx, client)
		if err != nil {
			fmt.Println("Mongo connection error while listing db's")
			fmt.Println(err)
			return err
		}

		for _, dbName := range dbs {
			db := client.Database(dbName)

			cols, err := mongo.ListCollections(ctx, db)
			if err != nil {
				fmt.Println("Mongo connection error")
				fmt.Println(err)
				return err
			}

			composer := render.NewFileComposer(render.TSOptions{
				RequiredThreshold: opts.OptionalThreshold,
				DateAs:            opts.DateAs,
				ObjectIDAs:        opts.ObjectIDAs,
				NullPolicy:        "optional",
				UseInterface:      false,
			})

			for _, colName := range cols {
				coll := db.Collection(colName)
				docs, err := mongo.SampleDocuments(ctx, coll, opts.Sample)
				if err != nil {
					return err
				}

				if len(docs) <= 0 {
					fmt.Printf("%s.%s -> sampled %d docs ( SKIPPING )\n", dbName, colName, len(docs))
					continue
				}

				schema := map[string]*infer.FieldStats{}
				totalDocs := 0

				for _, doc := range docs {
					infer.Flatten(doc, schema, &totalDocs)
				}

				tree := infer.BuildSchemaTree(schema)

				composer.AddCollection(tree, totalDocs, pkg.TypeNameFromCollection(colName))
			}
			fmt.Println(composer.String())
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&opts.URI, "uri", "", "MongoDB connection URI")
	rootCmd.Flags().StringVar(&opts.Out, "out", "", "Output TypeScript file path")
	rootCmd.Flags().IntVar(&opts.Sample, "sample", 2, "Sample size per collection")
	rootCmd.Flags().Float64Var(&opts.OptionalThreshold, "optional-threshold", 0.98, "Field required threshold based on samples")
	rootCmd.Flags().StringVar(&opts.DateAs, "date-as", "string", "How to emit dates: string|Date")
	rootCmd.Flags().StringVar(&opts.ObjectIDAs, "objectid-as", "string", "How to emit ObjectIds: string|ObjectId")
	rootCmd.Flags().StringVar(&opts.ConfigPath, "config", "", "Optional config path (yaml/json)")
	rootCmd.Flags().StringVar(&opts.EnvFile, "env-file", "", "Path to .env file (optional)")
}
