export function resolveTarget(platform, arch) {
  if (platform === "darwin" && arch === "arm64") return { binary: "bilink-darwin-arm64" };
  if (platform === "darwin" && arch === "x64") return { binary: "bilink-darwin-amd64" };
  if (platform === "linux" && arch === "x64") return { binary: "bilink-linux-amd64" };
  if (platform === "linux" && arch === "arm64") return { binary: "bilink-linux-arm64" };
  if (platform === "win32" && arch === "x64") return { binary: "bilink-windows-amd64.exe" };
  throw new Error(`unsupported platform: ${platform}-${arch}`);
}
