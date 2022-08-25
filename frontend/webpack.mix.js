let mix = require('laravel-mix');

mix.disableSuccessNotifications()
// mix.setPublicPath("./")
// mix.setResourceRoot("./")


mix.js('src/js/app.js', 'dist')
    .vue({
        version: 3
    })

mix.sass('src/scss/main.scss', 'dist')

mix.browserSync({
    proxy: {
        target: "http://app.loc:3000",
        ws: true,
    },
    open: 'http://app.loc:3000',
    browser: "google-chrome-stable",
    logConnections: true
});