# mongo-auto-type-gen
This is a usefull tool, that introspects your mongo collections and uses samples from them to generate types/interfaces/classes from them.

## Release
This repo publishes a single npm package that bundles prebuilt binaries for all platforms.

Release steps:
1. Update `package.json` version.
2. Tag and push, for example `v0.1.0`.
3. GitHub Actions builds binaries and runs `npm publish`.

Required secret:
- `NPM_TOKEN` with publish rights.
