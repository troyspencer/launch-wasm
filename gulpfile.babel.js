import gulp from 'gulp';
import browserSync from 'browser-sync';
import { exec } from 'child_process';

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
  serve
}

export default defaultTasks