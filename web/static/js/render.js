function setAuthMode(mode) {
  api.mode = mode;
  document.querySelectorAll('.segmented-btn').forEach((button) => {
    button.classList.toggle('active', button.dataset.mode === mode);
  });
  els.authSubmit.textContent = mode === 'register' ? 'Создать аккаунт' : 'Войти';
  els.authHint.textContent = mode === 'register'
    ? 'Регистрация создаёт пользователя и сразу выдаёт токен.'
    : 'Вход выдаёт JWT для favorites и защищённых запросов.';
}

function syncAuthUI() {
  if (api.token) {
    els.authBadge.textContent = 'authenticated';
    els.authBadge.className = 'status status-ok';
    els.authHint.textContent = 'JWT токен сохранён. Можно управлять избранным.';
    els.authGate.classList.add('is-hidden');
    document.body.classList.add('is-authenticated');
    window.scrollTo({ top: 0, behavior: 'instant' });
    return;
  }

  els.authBadge.textContent = 'guest';
  els.authBadge.className = 'status status-warn';
  els.authHint.textContent = 'Токен сохраняется в браузере и используется для favorites.';
  els.authGate.classList.remove('is-hidden');
  document.body.classList.remove('is-authenticated');
  window.scrollTo({ top: 0, behavior: 'instant' });
}

function updateStats() {
  els.booksCount.textContent = String(api.books.length);
  els.authorsCount.textContent = String(api.authors.length);
  els.categoriesCount.textContent = String(api.categories.length);
  els.favoritesCount.textContent = String(api.favorites.length);
}

function renderBooks() {
  if (!api.books.length) {
    els.booksList.innerHTML = '<div class="empty">Пока нет книг. Добавьте первую запись через форму слева.</div>';
    els.bookMeta.textContent = 'empty';
    return;
  }

  els.booksList.innerHTML = api.books.map((book) => `
    <article class="card">
      <div class="card-top">
        <div>
          <div class="card-title">${escapeHtml(book.title)}</div>
          <div class="card-meta">Автор: ${escapeHtml(book.author?.name || '—')} · Категория: ${escapeHtml(book.category?.name || '—')}</div>
        </div>
        <div class="pill">$${Number(book.price).toFixed(2)}</div>
      </div>
      <div class="card-actions">
        <button class="button button-ghost" data-action="edit-book" data-id="${book.id}">Редактировать</button>
        <button class="button button-ghost" data-action="delete-book" data-id="${book.id}">Удалить</button>
        <button class="button button-secondary" data-action="favorite-book" data-id="${book.id}">В избранное</button>
      </div>
    </article>
  `).join('');
  els.bookMeta.textContent = `${api.books.length} items`;
}

function renderFavorites() {
  if (!api.favorites.length) {
    els.favoritesList.innerHTML = '<div class="empty">Избранное пусто. Авторизуйтесь и добавьте книгу в favorites.</div>';
    return;
  }

  els.favoritesList.innerHTML = api.favorites.map((book) => `
    <article class="card">
      <div class="card-top">
        <div>
          <div class="card-title">${escapeHtml(book.title)}</div>
          <div class="card-meta">${escapeHtml(book.author?.name || '—')} · ${escapeHtml(book.category?.name || '—')}</div>
        </div>
        <div class="pill">$${Number(book.price).toFixed(2)}</div>
      </div>
      <div class="card-actions">
        <button class="button button-ghost" data-action="remove-favorite" data-id="${book.id}">Удалить из избранного</button>
      </div>
    </article>
  `).join('');
}

function fillBookForm(book) {
  els.bookForm.elements.id.value = book.id;
  els.bookForm.elements.title.value = book.title || '';
  els.bookForm.elements.author_id.value = book.author_id || book.author?.id || '';
  els.bookForm.elements.category_id.value = book.category_id || book.category?.id || '';
  els.bookForm.elements.price.value = book.price || '';
  els.bookFormMode.textContent = 'edit';
}

function clearBookForm() {
  els.bookForm.reset();
  els.bookForm.elements.id.value = '';
  els.bookFormMode.textContent = 'create';
}

function buildBookPayload(form) {
  return {
    title: form.elements.title.value.trim(),
    author_id: Number(form.elements.author_id.value),
    category_id: Number(form.elements.category_id.value),
    price: Number(form.elements.price.value),
  };
}

function buildBookFilters(form) {
  return {
    title: form.elements.title.value.trim(),
    author: form.elements.author.value,
    category: form.elements.category.value,
    min_price: form.elements.min_price.value,
    max_price: form.elements.max_price.value,
  };
}

function updateReferenceLists() {
  document.querySelectorAll('select[name="author"]').forEach((select) => populateSelect(select, api.authors, 'Все авторы'));
  document.querySelectorAll('select[name="category"]').forEach((select) => populateSelect(select, api.categories, 'Все категории'));
  document.querySelectorAll('select[name="author_id"]').forEach((select) => populateSelect(select, api.authors, 'Выберите автора'));
  document.querySelectorAll('select[name="category_id"]').forEach((select) => populateSelect(select, api.categories, 'Выберите категорию'));
}

function setActiveSidebarItem(view) {
  document.querySelectorAll('.notes-nav-item[data-view]').forEach((button) => {
    button.classList.toggle('active', button.dataset.view === view);
  });
  document.querySelectorAll('.view-chip').forEach((chip) => {
    chip.classList.toggle('active', chip.dataset.viewChip === view);
  });
}

function setWorkspaceView(view) {
  const targetView = view || 'home';
  api.currentView = targetView;
  setActiveSidebarItem(targetView);

  if (!api.token) {
    document.querySelectorAll('.dashboard-shell, [data-panel], .notes-mini-card').forEach((element) => {
      element.classList.add('is-hidden');
    });
    els.authGate.classList.remove('is-hidden');
    return;
  }

  const isHome = targetView === 'home';
  const shells = document.querySelectorAll('.dashboard-shell');
  shells.forEach((element) => {
    element.classList.toggle('is-hidden', !isHome && !element.matches(`[data-panel="${targetView}"]`));
  });

  const allPanels = document.querySelectorAll('[data-panel]');
  allPanels.forEach((element) => {
    const visible = isHome || element.dataset.panel === targetView;
    element.classList.toggle('is-hidden', !visible);
  });

  document.querySelectorAll('.notes-mini-card').forEach((element) => {
    const visible = isHome || targetView === 'reminders' || targetView === 'create-label';
    element.classList.toggle('is-hidden', !visible);
  });
}

function showAllWorkspaceSections() {
  setWorkspaceView(api.currentView || 'home');
}