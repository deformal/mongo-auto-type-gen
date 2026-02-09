package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deformal/mongo-auto-type-gen/core/internals/config"
	"github.com/deformal/mongo-auto-type-gen/core/internals/infer"
	"github.com/deformal/mongo-auto-type-gen/core/internals/mongo"
	"github.com/deformal/mongo-auto-type-gen/core/internals/render"
	"github.com/deformal/mongo-auto-type-gen/core/pkg"
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

		if strings.TrimSpace(opts.Out) == "" {
			return fmt.Errorf("--out flag is required")
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

			out := composer.String()

			if strings.TrimSpace(out) == "" {
				continue
			}

			outPath := opts.Out

			if len(dbs) > 1 {
				outPath = filepath.Join(opts.Out, dbName+".ts")
			} else {
				if strings.HasSuffix(outPath, string(os.PathSeparator)) {
					outPath = filepath.Join(outPath, dbName+".ts")
				} else if info, err := os.Stat(outPath); err == nil && info.IsDir() {
					outPath = filepath.Join(outPath, dbName+".ts")
				} else if filepath.Ext(outPath) == "" {
					outPath = outPath + ".ts"
				}
			}

			if err := render.WriteFile(outPath, out); err != nil {
				return err
			}

			fmt.Printf("wrote %s\n", outPath)
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
