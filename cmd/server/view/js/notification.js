(() => {
  let timer = undefined;
  function hideSnackbar(e) {
    if (e.detail.type !== 'notification') {
      return;
    }
    if (timer) {
      clearTimeout(timer);
    }

    timer = setTimeout(() => {
      const snackbar = document.getElementById("snackbar");
      snackbar.classList.add("hidden");
    }, 3000);
  }

  document.body.addEventListener('htmx:sseMessage', hideSnackbar);
  document.body.addEventListener('htmx:sseClose', function() {
    document.body.removeEventListener('htmx:sseBeforeMessage', hideSnackbar);
  });
})()
