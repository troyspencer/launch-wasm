import gulp from 'gulp';
import browserSync from 'browser-sync';
import { exec } from 'child_process';
import workboxBuild from 'workbox-build';

const browserSyncInstance = browserSync.create()

const buildWasm = (cb) => {
  exec('env GOOS=js GOARCH=wasm go build -o server/main.wasm go/main.go', (err, stdout, stderr) => {
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.log(stderr);
    }      
    cb(err)
    reloadServer()
  })
}

const dockerBuildWasm = (cb) => {
  exec('rm -f server/static/main.wasm.gz && rm -f server/static/main.wasm && env GOOS=js GOARCH=wasm go build -o server/static/main.wasm go/main.go && gzip -k server/static/main.wasm', (err, stdout, stderr) => {
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.log(stderr);
    }
    cb(err)
  })
}

// NOTE: This should be run *AFTER* all your assets are built
const buildSW = () => {
  // This will return a Promise
  return workboxBuild.generateSW({
    globDirectory: './server/static',
    globPatterns: [
      '**\/*.{html,js,json,}',
    ],
    swDest: './server/static/sw.js',
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
  var folders = ["","contact/","world/"]
  for (var i = 0; i < folders.length; i++) {
    gulp.watch("./go/"+folders[i]+"*.go", dockerBuildWasm)
  }
}


const reloadServer = () => {
  setTimeout(() => { browserSyncInstance.reload() }, 300); 
}

const serve = () => {
  
  browserSyncInstance.init({
    "callbacks": {
      ready: (err, bs) => {
        bs.utils.serveStatic.mime.define({ 'application/wasm': ['wasm'] });
      }
    },
    server: {
        baseDir: "./server"
    },
    "browser": "google chrome",
    "open": false
  });

  var folders = ["","contact/","world/"]

  for (var i = 0; i < folders.length; i++) {
    gulp.watch("./go/"+folders[i]+"*.go", buildWasm)
  }
}

const defaultTasks = gulp.series(serve)

export {
  buildWasm,
  reloadServer,
  serve,
  watch,
  dockerBuildWasm,
  buildSW
}

export default defaultTasks
