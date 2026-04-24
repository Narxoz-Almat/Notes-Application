async function loadReferenceData() {
  const [authors, categories] = await Promise.all([
    request('/authors'),
    request('/categories'),
  ]);

  api.authors = authors.data || [];
  api.categories = categories.data || [];
  updateReferenceLists();
  updateStats();
}

async function loadBooks() {
  const params = new URLSearchParams();
  params.set('page', String(api.page));
  params.set('limit', String(api.limit));

  Object.entries(api.filters).forEach(([key, value]) => {
    if (value) params.set(key, value);
  });

  const data = await request(`/books?${params.toString()}`);
  api.books = data.data || [];
  els.pageInfo.textContent = `page ${data.page || api.page} / total ${data.total || 0}`;
  renderBooks();
  updateStats();
}

async function loadFavorites() {
  if (!api.token) {
    api.favorites = [];
    renderFavorites();
    updateStats();
    return;
  }

  try {
    const data = await request('/books/favorites?page=1&limit=12');
    api.favorites = data.data || [];
  } catch (error) {
    api.favorites = [];
  }

  renderFavorites();
  updateStats();
}

async function refreshAll() {
  try {
    await loadReferenceData();
    await loadBooks();
    await loadFavorites();
  } catch (error) {
    showToast(error.message);
  }
}