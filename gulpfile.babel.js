import gulp from 'gulp';
import { exec } from 'child_process';

const buildWasm = (cb) => {
  exec('rm -f react/static/main.wasm.gz && rm -f react/static/main.wasm && env GOOS=js GOARCH=wasm go build -o react/static/main.wasm game/main.go && gzip -k react/static/main.wasm', (err, stdout, stderr) => {
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.log(stderr);
    }
    cb(err)
  })
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
  watch
}

export default defaultTasks
