// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

const site = "https://halftoothed.github.io/gostman";

// https://astro.build/config
export default defineConfig({
  site: "https://halftoothed.github.io",
  base: "/gostman",
  outDir: './dist',
  integrations: [
    starlight({
      title: "Gostman",
      head: [
        {
          tag: "meta",
          attrs: { property: "og:image", content: site + "/og.jpg?v=1" },
        },
        {
          tag: "meta",
          attrs: { property: "twitter:image", content: site + "/og.jpg?v=1" },
        },
        {
          tag: "link",
          attrs: { rel: "preconnect", href: "https://fonts.googleapis.com" },
        },
        {
          tag: "link",
          attrs: {
            rel: "preconnect",
            href: "https://fonts.gstatic.com",
            crossorigin: true,
          },
        },
        {
          tag: "link",
          attrs: {
            rel: "stylesheet",
            href: "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@500;600&display=swap",
          },
        },
        {
          tag: "script",
          attrs: {
            src: "https://cdn.jsdelivr.net/npm/@minimal-analytics/ga4/dist/index.js",
            async: true,
          },
        },
        {
          tag: "script",
          content: ` window.minimalAnalytics = {
            trackingId: 'G-WFLBCRZ7MC',
            autoTrack: true,
          };`,
        },
      ],
      social: {
        github: "https://github.com/HalfToothed/gostman",
      },
      customCss: ["./src/styles/custom.css"],
    }),
  ],
});
