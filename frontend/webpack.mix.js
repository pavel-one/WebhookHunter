let mix = require('laravel-mix');

mix.disableSuccessNotifications()
// mix.setPublicPath("./")
// mix.setResourceRoot("./")

mix.js('src/js/app.js', 'dist')
    .vue({
        version: 3
    })
mix.sass('src/scss/main.scss', 'dist')