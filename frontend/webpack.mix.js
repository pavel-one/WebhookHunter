let mix = require('laravel-mix');

mix.js('src/js/app.js', 'dist')
    .setPublicPath('dist')
    .vue({
        version: 3
    })

mix.sass('src/scss/main.scss', 'dist')
    .setPublicPath('dist')