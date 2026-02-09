# mongo-auto-type-gen
This is a usefull tool, that introspects your mongo collections and uses samples from them to generate types/interfaces/classes from them.

## Usage
Generate TypeScript types from MongoDB collections by inference.

### Command
`mongots`

### Examples
- `mongots --uri mongodb://localhost:27017 --out ./generated`
- `mongots --env-file .env --out ./generated`

### Flags
- `--uri <string>` MongoDB connection URI
- `--out <path>` Output TypeScript file path (required)
- `--sample <int>` Sample size per collection (default: 2)
- `--optional-threshold <float>` Field required threshold (default: 0.98)
- `--date-as <string>` `string|Date` (default: `string`)
- `--objectid-as <string>` `string|ObjectId` (default: `string`)
- `--config <path>` Optional config path (yaml/json)
- `--env-file <path>` Path to .env file (optional)

### Environment variables
- `MONGOTS_MONGO_URI` or `MONGO_URI`
- `MONGOTS_OUT`
- `MONGOTS_SAMPLE`
- `MONGOTS_OPTIONAL_THRESHOLD`
- `MONGOTS_DATE_AS`
- `MONGOTS_OBJECTID_AS`
