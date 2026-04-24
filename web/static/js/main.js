async function initApp() {
  syncAuthUI();
  setAuthMode(api.mode);
  bindEvents();
  setWorkspaceView('home');
  window.scrollTo({ top: 0, behavior: 'instant' });
  await refreshAll();
}

initApp();