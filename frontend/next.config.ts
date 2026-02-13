import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Use standalone output for Docker deployment (without static export)
  output: "standalone",
};

export default nextConfig;
