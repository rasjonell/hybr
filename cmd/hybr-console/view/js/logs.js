function clearLogs() {
  const logs = document.getElementById("log-list");
  logs.innerHTML = "";
}

document.body.addEventListener('htmx:sseMessage', function(e) {
  if (e.detail.type !== "log") {
    return;
  }

  const logs = document.getElementById("log-list");
  logs.scrollTo(0, logs.scrollHeight)
})
