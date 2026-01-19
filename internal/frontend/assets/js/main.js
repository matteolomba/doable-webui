function doableApp() {
    return {
        todos: [],
        lists: [],
        searchQuery: '',
        searchField: '3', // 1: Title, 2: Description, 3: Both (Default)
        statusFilter: '1', // 1: Incomplete, 2: All, 3: Complete
        isLoading: true,

        // Modal State
        isModalOpen: false,
        isChangelogOpen: false,
        currentTodoId: null,
        formData: {
            title: '',
            description: '',
            listId: ''
        },

        // List Manager State
        activeListFilter: null,
        isListModalOpen: false,
        isDeleteModalOpen: false,
        editingListId: null,
        deleteListId: null,
        listFormData: {
            name: '',
            color: '#4f46e5',
            icon: ''
        },

        // Configuration
        presetColors: [
            '#ef4444', '#f97316', '#f59e0b', '#84cc16',
            '#10b981', '#06b6d4', '#0ea5e9', '#3b82f6',
            '#6366f1', '#8b5cf6', '#d946ef', '#f43f5e'
        ],
        availableIcons: [
            // General
            { name: 'list' },
            { name: 'home' },
            { name: 'star' },
            { name: 'heart' },

            // Work & Dev
            { name: 'briefcase' },
            { name: 'code' },
            { name: 'terminal' },
            { name: 'server' },
            { name: 'cpu' },

            // Tracking & Shipping
            { name: 'package' },
            { name: 'map-pin' },

            // DIY & Repair
            { name: 'tool' },
            { name: 'settings' },

            // Media
            { name: 'image' },
            { name: 'camera' },
            { name: 'instagram' },
            { name: 'video' },

            // Social & Events
            { name: 'calendar' },
            { name: 'users' },
            { name: 'smile' },

            // Shopping & Money
            { name: 'shopping-cart' },
            { name: 'tag' },
            { name: 'dollar-sign' },

            // School & Learning
            { name: 'book' },
            { name: 'bookmark' },

            // Lifestyle & Others
            { name: 'coffee' },
            { name: 'globe' },
            { name: 'gift' },
            { name: 'zap' },
            { name: 'music' }
        ],

        // Notification System
        notification: {
            show: false,
            message: '',
            type: 'success'
        },

        // Theme State
        isDark: false,
        mobileView: 'todos', // 'todos' or 'lists'

        async init() {
            // Theme Init
            if (localStorage.getItem('theme') === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
                this.isDark = true;
                document.documentElement.classList.add('dark');
            } else {
                this.isDark = false;
                document.documentElement.classList.remove('dark');
            }

            await this.fetchData();
        },

        getIcon(iconName, width = 24, height = 24, className = '') {
            if (typeof feather !== 'undefined' && feather.icons[iconName]) {
                return feather.icons[iconName].toSvg({ width, height, class: className });
            }
            return '';
        },

        toggleTheme() {
            this.isDark = !this.isDark;
            if (this.isDark) {
                document.documentElement.classList.add('dark');
                localStorage.setItem('theme', 'dark');
            } else {
                document.documentElement.classList.remove('dark');
                localStorage.setItem('theme', 'light');
            }
        },

        // List Filter Logic
        toggleListFilter(listId) {
            if (this.activeListFilter === listId) {
                this.activeListFilter = null; // Toggle off
            } else {
                this.activeListFilter = listId;
            }
        },

        async fetchData() {
            this.isLoading = true;
            try {
                // Parallel fetch for speed
                const [listsRes, todosRes] = await Promise.all([
                    fetch('api/lists'),
                    fetch('api/todos')
                ]);

                if (listsRes.ok) this.lists = await listsRes.json() || [];
                if (todosRes.ok) this.todos = await todosRes.json() || [];

                // Sort todos by last modified
                this.todos.sort((a, b) => new Date(b.lastModified) - new Date(a.lastModified));
            } catch (error) {
                console.error("Fetch error:", error);
                this.showNotification("Errore di connessione", "error");
            } finally {
                this.isLoading = false;
            }
        },

        get filteredTodos() {
            return this.todos.filter(todo => {
                // List Filter
                if (this.activeListFilter) {
                    if (this.activeListFilter === 'general') {
                        if (todo.listId && todo.listId !== '') return false;
                    } else {
                        if (todo.listId !== this.activeListFilter) return false;
                    }
                }

                // Status Filter
                if (this.statusFilter === '1' && todo.isCompleted) return false;
                if (this.statusFilter === '3' && !todo.isCompleted) return false;

                // Search Filter
                if (this.searchQuery) {
                    const q = this.searchQuery.toLowerCase();
                    const titleMatch = todo.title.toLowerCase().includes(q);
                    const descMatch = todo.description ? todo.description.toLowerCase().includes(q) : false;

                    if (this.searchField === '1' && !titleMatch) return false;
                    if (this.searchField === '2' && !descMatch) return false;
                    if (this.searchField === '3' && !titleMatch && !descMatch) return false;
                }
                return true;
            });
        },

        getList(listId) {
            return this.lists.find(l => l.id === listId);
        },

        // ================= Todos Actions =================
        async toggleTodo(id) {
            const todo = this.todos.find(t => t.id === id);
            if (!todo) return;

            // Optimistic Update
            const originalState = todo.isCompleted;
            todo.isCompleted = !originalState;

            try {
                const endpoint = todo.isCompleted ? `api/todos/${id}/check` : `api/todos/${id}/uncheck`;
                const res = await fetch(endpoint, { method: 'PUT' });

                if (!res.ok) throw new Error('Failed to update');

                // Refresh data in background to ensure consistency
                // await this.fetchData(); 
            } catch (error) {
                todo.isCompleted = originalState; // Revert
                this.showNotification("Errore aggiornamento", "error");
            }
        },

        openAddModal() {
            this.currentTodoId = null;
            this.formData = { title: '', description: '', listId: this.activeListFilter && this.activeListFilter !== 'general' ? this.activeListFilter : '' };
            this.isModalOpen = true;
        },

        editTodo(id) {
            const todo = this.todos.find(t => t.id === id);
            if (!todo) return;

            this.currentTodoId = id;
            this.formData = {
                title: todo.title,
                description: todo.description,
                listId: todo.listId || ''
            };
            this.isModalOpen = true;
        },

        closeModal() {
            this.isModalOpen = false;
        },

        async saveTodo() {
            try {
                const method = this.currentTodoId ? 'PUT' : 'POST';
                const url = this.currentTodoId ? `api/todos/${this.currentTodoId}` : 'api/todos';

                const res = await fetch(url, {
                    method,
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(this.formData)
                });

                if (res.ok) {
                    const data = await res.json();

                    if (this.currentTodoId) {
                        // Update existing
                        const idx = this.todos.findIndex(t => t.id === this.currentTodoId);
                        if (idx !== -1) this.todos[idx] = { ...this.todos[idx], ...data }; // Merge to keep other fields
                        this.showNotification("Attività aggiornata");
                    } else {
                        // Add new
                        this.todos.unshift(data);
                        this.showNotification("Attività creata");
                    }
                    this.closeModal();
                } else {
                    throw new Error("API Error");
                }
            } catch (error) {
                console.error(error);
                this.showNotification("Errore nel salvataggio", "error");
            }
        },

        async deleteTodoConfirm(id) {
            if (!confirm('Sei sicuro di voler eliminare questa attività?')) return;

            try {
                const res = await fetch(`api/todos/${id}`, { method: 'DELETE' });
                if (res.ok) {
                    this.todos = this.todos.filter(t => t.id !== id);
                    this.showNotification("Attività eliminata");
                } else {
                    throw new Error("Delete failed");
                }
            } catch (error) {
                this.showNotification("Errore eliminazione", "error");
            }
        },

        // ================= Lists Actions =================
        openCreateListModal() {
            this.editingListId = null;
            this.listFormData = { name: '', color: this.presetColors[4], icon: 'list' }; // Default color emerald-500, icon list
            this.isListModalOpen = true;
        },

        openEditListModal(listId) {
            const list = this.getList(listId);
            if (!list) return;
            this.editingListId = listId;
            const hexColor = this.rgbToHex(list.color);
            this.listFormData = { name: list.name, color: hexColor, icon: list.icon || 'list' };
            this.isListModalOpen = true;
        },

        closeListModal() {
            this.isListModalOpen = false;
        },

        openChangelogModal() {
            this.isChangelogOpen = true;
        },

        closeChangelogModal() {
            this.isChangelogOpen = false;
        },

        async saveList() {
            if (!this.listFormData.name.trim()) return;

            try {
                const rgb = this.hexToRgb(this.listFormData.color);
                const body = JSON.stringify({
                    name: this.listFormData.name,
                    color: rgb,
                    icon: this.listFormData.icon
                });

                let res;
                if (this.editingListId) {
                    // Edit
                    res = await fetch(`api/lists/${this.editingListId}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: body
                    });
                } else {
                    // Create
                    res = await fetch('api/lists', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: body
                    });
                }

                if (res.ok) {
                    const data = await res.json();
                    if (this.editingListId) {
                        const idx = this.lists.findIndex(l => l.id === this.editingListId);
                        if (idx !== -1) this.lists[idx] = data;
                        this.showNotification("Lista aggiornata");
                    } else {
                        this.lists.push(data);
                        this.showNotification("Lista creata");
                    }
                    this.closeListModal();
                } else {
                    throw new Error("List save failed");
                }
            } catch (error) {
                console.error(error);
                this.showNotification("Errore salvataggio lista", "error");
            }
        },

        promptDeleteList(id) {
            this.deleteListId = id;
            this.isDeleteModalOpen = true;
        },

        async confirmDeleteList() {
            if (!this.deleteListId) return;

            try {
                const res = await fetch(`api/lists/${this.deleteListId}`, { method: 'DELETE' });
                if (res.ok) {
                    this.lists = this.lists.filter(l => l.id !== this.deleteListId);
                    // If we were filtering by this list, reset filter
                    if (this.activeListFilter === this.deleteListId) {
                        this.activeListFilter = null;
                    }
                    this.showNotification("Lista eliminata");
                } else {
                    throw new Error("Delete list failed");
                }
            } catch (error) {
                this.showNotification("Errore eliminazione lista", "error");
            } finally {
                this.isDeleteModalOpen = false;
                this.deleteListId = null;
            }
        },

        // ================= Helpers =================
        showNotification(message, type = 'success') {
            this.notification = { show: true, message, type };
            setTimeout(() => {
                this.notification.show = false;
            }, 3000);
        },

        hexToRgb(hex) {
            const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
            return result ? [
                parseInt(result[1], 16),
                parseInt(result[2], 16),
                parseInt(result[3], 16)
            ] : [0, 0, 0];
        },

        rgbToHex(rgb) {
            if (!rgb || rgb.length < 3) return '#000000';
            return "#" + ((1 << 24) + (rgb[0] << 16) + (rgb[1] << 8) + rgb[2]).toString(16).slice(1);
        },

        getContrastColor(rgb) {
            if (!rgb) return '#000000';
            // Simple luminance formula
            const yiq = ((rgb[0] * 299) + (rgb[1] * 587) + (rgb[2] * 114)) / 1000;
            return (yiq >= 128) ? '#000000' : '#ffffff';
        },

        // Darkens a color by a percentage (0-100) strictly for text visibility
        darkenColor(rgb, percent) {
            if (!rgb) return 'rgb(0,0,0)';

            // Check if it's already dark enough
            const yiq = ((rgb[0] * 299) + (rgb[1] * 587) + (rgb[2] * 114)) / 1000;
            if (yiq < 50) {
                // It's already dark, return as is or slightly brighter?
                // Let's just return it as rgb string
                return `rgb(${rgb.join(',')})`;
            }

            const f = (val) => Math.max(0, Math.floor(val * (100 - percent) / 100));
            return `rgb(${f(rgb[0])},${f(rgb[1])},${f(rgb[2])})`;
        },

        // For dark mode: lightens color
        lightenColor(rgb, percent) {
            if (!rgb) return 'rgb(255,255,255)';
            const f = (val) => Math.min(255, Math.floor(val + (255 - val) * percent / 100));
            return `rgb(${f(rgb[0])},${f(rgb[1])},${f(rgb[2])})`;
        },

        // Smart text color for chips:
        // Use strict White/Black based on luminance of the list color
        getChipTextColor(listId) {
            const list = this.getList(listId);
            // Default fallback
            if (!list) return this.isDark ? '#e5e7eb' : '#1f2937';

            return this.getContrastColor(list.color);
        },

        getChipBgColor(listId) {
            const list = this.getList(listId);
            // Default fallback
            if (!list) return this.isDark ? 'rgba(75, 85, 99, 0.4)' : '#e5e7eb';

            // Solid color for chips as requested
            return `rgb(${list.color.join(',')})`;
        },

        // New: Card background tint based on list
        getTodoBgColor(listId) {
            const list = this.getList(listId);
            // User requested high contrast/separation. 
            // We strip the tint to make cards "pop" as white/dark blocks against the page.
            if (this.isDark) {
                return '#1e293b'; // Slate 800 (Surface Variant)
            } else {
                return '#ffffff'; // White
            }
        },

        // Helper: Border color for cards to match list (More visible now)
        getTodoBorderColor(listId) {
            const list = this.getList(listId);
            if (!list) return this.isDark ? '#475569' : '#cbd5e1'; // Slate 600 : Slate 300

            // Stronger border opacity
            return `rgba(${list.color.join(',')}, ${this.isDark ? 0.4 : 0.5})`;
        }
    }
}