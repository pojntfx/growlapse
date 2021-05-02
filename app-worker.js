const cacheName = "app-" + "e5b2d1a6dcf51b686a6f190b460db1425cae4656";

self.addEventListener("install", event => {
  console.log("installing app worker e5b2d1a6dcf51b686a6f190b460db1425cae4656");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "",
        "/growlapse",
        "/growlapse/app.css",
        "/growlapse/app.js",
        "/growlapse/manifest.webmanifest",
        "/growlapse/wasm_exec.js",
        "/growlapse/web/app.wasm",
        "/growlapse/web/icon.png",
        "/growlapse/web/index.css",
        "https://unpkg.com/@patternfly/patternfly@4.96.2/patternfly-addons.css",
        "https://unpkg.com/@patternfly/patternfly@4.96.2/patternfly.css",
        
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker e5b2d1a6dcf51b686a6f190b460db1425cae4656 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
