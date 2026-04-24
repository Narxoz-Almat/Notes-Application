async function handleAuthSubmit(event) {
  event.preventDefault();
  const form = event.currentTarget;
  const payload = {
    name: form.elements.name.value.trim(),
    email: form.elements.email.value.trim(),
    password: form.elements.password.value,
  };

  try {
    const path = api.mode === 'register' ? '/auth/register' : '/auth/login';
    const result = await request(path, { method: 'POST', body: JSON.stringify(payload) });
    api.token = result.token;
    localStorage.setItem('bookstore-token', api.token);
    syncAuthUI();
    setWorkspaceView(api.currentView || 'home');
    showToast(api.mode === 'register' ? 'Пользователь создан' : 'Успешный вход');
    await loadFavorites();
  } catch (error) {
    showToast(error.message);
  }
}

async function handleLogout() {
  api.token = '';
  localStorage.removeItem('bookstore-token');
  syncAuthUI();
  setWorkspaceView('home');
  await loadFavorites();
  showToast('Сессия очищена');
}

async function handleAuthorSubmit(event) {
  event.preventDefault();
  try {
    await request('/authors', {
      method: 'POST',
      body: JSON.stringify({ name: event.currentTarget.elements.name.value.trim() }),
    });
    event.currentTarget.reset();
    showToast('Автор добавлен');
    await refreshAll();
  } catch (error) {
    showToast(error.message);
  }
}

async function handleCategorySubmit(event) {
  event.preventDefault();
  try {
    await request('/categories', {
      method: 'POST',
      body: JSON.stringify({ name: event.currentTarget.elements.name.value.trim() }),
    });
    event.currentTarget.reset();
    showToast('Категория добавлена');
    await refreshAll();
  } catch (error) {
    showToast(error.message);
  }
}

async function handleBookSubmit(event) {
  event.preventDefault();
  const form = event.currentTarget;
  const id = form.elements.id.value;
  const payload = buildBookPayload(form);

  try {
    if (id) {
      await request(`/books/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
      showToast('Книга обновлена');
    } else {
      await request('/books', { method: 'POST', body: JSON.stringify(payload) });
      showToast('Книга создана');
    }
    clearBookForm();
    await refreshAll();
  } catch (error) {
    showToast(error.message);
  }
}

async function handleFilterSubmit(event) {
  event.preventDefault();
  api.page = 1;
  api.filters = buildBookFilters(event.currentTarget);
  await loadBooks();
}

async function handleResetFilters() {
  els.filterForm.reset();
  api.page = 1;
  api.filters = {};
  await loadBooks();
}

async function handlePrevPage() {
  if (api.page === 1) return;
  api.page -= 1;
  await loadBooks();
}

async function handleNextPage() {
  api.page += 1;
  await loadBooks();
}

async function handleBooksListClick(event) {
  const action = event.target.dataset.action;
  const id = event.target.dataset.id;
  if (!action || !id) return;

  try {
    if (action === 'edit-book') {
      const book = api.books.find((item) => String(item.id) === String(id));
      if (book) fillBookForm(book);
      return;
    }

    if (action === 'delete-book') {
      await request(`/books/${id}`, { method: 'DELETE' });
      showToast('Книга удалена');
      await refreshAll();
      return;
    }

    if (action === 'favorite-book') {
      await request(`/books/${id}/favorites`, { method: 'PUT' });
      showToast('Добавлено в избранное');
      await loadFavorites();
    }
  } catch (error) {
    showToast(error.message);
  }
}

async function handleFavoritesClick(event) {
  const action = event.target.dataset.action;
  const id = event.target.dataset.id;
  if (action !== 'remove-favorite' || !id) return;

  try {
    await request(`/books/${id}/favorites`, { method: 'DELETE' });
    showToast('Удалено из избранного');
    await loadFavorites();
  } catch (error) {
    showToast(error.message);
  }
}

function bindEvents() {
  els.authForm.addEventListener('submit', handleAuthSubmit);
  els.logoutBtn.addEventListener('click', handleLogout);
  els.authorForm.addEventListener('submit', handleAuthorSubmit);
  els.categoryForm.addEventListener('submit', handleCategorySubmit);
  els.bookForm.addEventListener('submit', handleBookSubmit);
  els.clearBookForm.addEventListener('click', clearBookForm);
  els.filterForm.addEventListener('submit', handleFilterSubmit);
  els.resetFilters.addEventListener('click', handleResetFilters);
  els.prevPage.addEventListener('click', handlePrevPage);
  els.nextPage.addEventListener('click', handleNextPage);
  els.booksList.addEventListener('click', handleBooksListClick);
  els.favoritesList.addEventListener('click', handleFavoritesClick);

  els.sidebarNav.addEventListener('click', (event) => {
    const button = event.target.closest('[data-view]');
    if (!button) return;
    setWorkspaceView(button.dataset.view);
  });

  document.querySelectorAll('[data-view-chip]').forEach((chip) => {
    chip.addEventListener('click', () => setWorkspaceView(chip.dataset.viewChip));
  });

  document.querySelectorAll('.segmented-btn').forEach((button) => {
    button.addEventListener('click', () => setAuthMode(button.dataset.mode));
  });
}