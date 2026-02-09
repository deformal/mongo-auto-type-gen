#!/usr/bin/env node
"use strict";

const { spawnSync } = require("node:child_process");
const fs = require("node:fs");
const path = require("node:path");

const platform = process.platform;
const arch = (() => {
  if (process.arch === "x64") return "amd64";
  if (process.arch === "arm64") return "arm64";
  return process.arch;
})();

const binName = (() => {
  if (platform === "win32") return "mongots-win32-" + arch + ".exe";
  if (platform === "darwin") return "mongots-darwin-" + arch;
  return "mongots-" + platform + "-" + arch;
})();

const distDir = path.join(__dirname, "..", "dist");
const binPath = path.join(distDir, binName);
const fallbackPath = path.join(distDir, "mongots");

const resolved = fs.existsSync(binPath) ? binPath : fallbackPath;

if (!fs.existsSync(resolved)) {
  console.error("mongots binary not found.");
  console.error("Tried: " + binPath);
  console.error("Also tried: " + fallbackPath);
  console.error("If you installed from source, run: npm run build:go");
  process.exit(1);
}

const result = spawnSync(resolved, process.argv.slice(2), {
  stdio: "inherit",
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 1);
