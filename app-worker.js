const cacheName = "app-" + "e7f9e5b253a0b438207463e12bce1e7f3ac40702";

self.addEventListener("install", event => {
  console.log("installing app worker e7f9e5b253a0b438207463e12bce1e7f3ac40702");
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
  console.log("app worker e7f9e5b253a0b438207463e12bce1e7f3ac40702 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
