(() => {
  const defaultContainerClass = "flex flex-row w-fit items-center px-2 py-1 rounded-lg capitalize"
  const successContainerClass = `${defaultContainerClass} bg-green-200 text-green-800`;
  const errorContainerClass = `${defaultContainerClass} bg-red-200 text-red-800`;

  const defaultChipOuterClass = "absolute inline-flex h-full w-full animate-ping rounded-full opacity-75";
  const successChipOuterClass = `${defaultChipOuterClass} bg-green-400`;
  const errorChipOuterClass = `${defaultChipOuterClass} bg-red-400`;

  const defaultChipInnerClass = "relative inline-flex size-3 rounded-full";
  const successChipInnerClass = `${defaultChipInnerClass} bg-green-500`;
  const errorChipInnerClass = `${defaultChipInnerClass} bg-red-500`;

  function updateStyles(e) {
    if (!e.detail.type.startsWith("status")) {
      return
    }

    if (e.detail.elt.innerText.toLowerCase() === "running") {
      e.detail.elt.parentElement.className = successContainerClass;
      e.detail.elt.previousElementSibling.childNodes[0].className = successChipOuterClass;
      e.detail.elt.previousElementSibling.childNodes[2].className = successChipInnerClass;
    } else {
      e.detail.elt.parentElement.className = errorContainerClass;
      e.detail.elt.previousElementSibling.childNodes[0].className = errorChipOuterClass;
      e.detail.elt.previousElementSibling.childNodes[2].className = errorChipInnerClass;
    }
  }

  document.body.addEventListener('htmx:sseMessage', updateStyles);
  document.body.addEventListener('htmx:sseClose', function() {
    document.body.removeEventListener('htmx:sseBeforeMessage', updateStyles);
  });
})()
