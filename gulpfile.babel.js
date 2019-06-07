import gulp from 'gulp';
import { exec } from 'child_process';
import workboxBuild from 'workbox-build';

const buildWasm = (cb) => {
  exec('rm -f react/static/main.wasm.gz && rm -f react/static/main.wasm && env GOOS=js GOARCH=wasm go build -o react/static/main.wasm game/main.go && gzip -k react/static/main.wasm', (err, stdout, stderr) => {
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.log(stderr);
    }
    cb(err)
    buildSW()
  })
}

// NOTE: This should be run *AFTER* all your assets are built
const buildSW = () => {
  // This will return a Promise
  return workboxBuild.generateSW({
    globDirectory: './react/static',
    globPatterns: [
      '**\/*.{html,js,json,}',
    ],
    swDest: './react/static/sw.js',
    // Define runtime caching rules.
    runtimeCaching: [{
      // Match any request ends with .png, .jpg, .jpeg or .svg.
      urlPattern: /\.(?:png|jpg|jpeg|svg)$/,

      // Apply a cache-first strategy.
      handler: 'StaleWhileRevalidate',

      options: {
        // Use a custom cache name.
        cacheName: 'images',

        // Only cache 10 images.
        expiration: {
          maxEntries: 10,
        },
      },
    },{
      // Match any request ends with .png, .jpg, .jpeg or .svg.
      urlPattern: /\.(?:wasm)$/,

      // Apply a cache-first strategy.
      handler: 'StaleWhileRevalidate',

      options: {
        // Use a custom cache name.
        cacheName: 'refresh',

        expiration: {
          maxEntries: 10,
        },
      },
    }],
  });
}

const watch = () => {
  var folders = ["","/contact","/world","/bodies"]
  for (var i = 0; i < folders.length; i++) {
    gulp.watch("./game"+folders[i]+"/*.go", buildWasm)
  }
}

const defaultTasks = gulp.series(watch)

export {
  buildWasm,
  buildSW,
  watch
}

export default defaultTasks
