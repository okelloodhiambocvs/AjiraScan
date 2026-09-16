function startLoading() {
  const button = document.getElementById("analyzeBtn");
  const loadingBar = document.getElementById("loadingBar");
  if (button) {
    button.innerText = "Analyzing";
    button.disabled = true;
  }
  if (loadingBar) {
    loadingBar.classList.remove("hidden");
  }
}
