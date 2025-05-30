// Current path stack
let currentPath = [];
let fileSystem = null;

// DOM Elements
const fileExplorer = document.getElementById('file-explorer');
const fileContent = document.getElementById('file-content');
const fileTitle = document.getElementById('file-title');
const breadcrumb = document.getElementById('breadcrumb');
const refreshBtn = document.getElementById('refresh-btn');
const copyBtn = document.getElementById('copy-btn');
const downloadBtn = document.getElementById('download-btn');
const fileActions = document.getElementById('file-actions');

// File type icons
const FILE_ICONS = {
    'directory': 'folder',
    'file': 'file-earmark',
    'json': 'filetype-json',
    'js': 'filetype-js',
    'css': 'filetype-css',
    'html': 'filetype-html',
    'md': 'file-earmark-text',
    'py': 'filetype-py',
    'jpg': 'file-image',
    'jpeg': 'file-image',
    'png': 'file-image',
    'gif': 'file-image',
    'svg': 'file-image'
};

// Get file icon based on file extension
function getFileIcon(name, type) {
    if (type === 'directory') return 'folder';
    const ext = name.split('.').pop().toLowerCase();
    return FILE_ICONS[ext] || 'file-earmark';
}

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Render breadcrumb
function renderBreadcrumb() {
    breadcrumb.innerHTML = '';
    
    // Add root/home
    const homeItem = document.createElement('li');
    homeItem.className = 'breadcrumb-item' + (currentPath.length === 0 ? ' active' : '');
    const homeLink = document.createElement('a');
    homeLink.href = '#';
    homeLink.textContent = 'Home';
    homeLink.addEventListener('click', (e) => {
        e.preventDefault();
        navigateTo([]);
    });
    homeItem.appendChild(homeLink);
    breadcrumb.appendChild(homeItem);
    
    // Add path segments
    let currentPathStr = '';
    currentPath.forEach((segment, index) => {
        currentPathStr += (currentPathStr ? '/' : '') + segment;
        const isLast = index === currentPath.length - 1;
        
        const item = document.createElement('li');
        item.className = 'breadcrumb-item' + (isLast ? ' active' : '');
        
        if (!isLast) {
            const link = document.createElement('a');
            link.href = '#';
            link.textContent = segment;
            const pathToNavigate = currentPath.slice(0, index + 1);
            link.addEventListener('click', (e) => {
                e.preventDefault();
                navigateTo(pathToNavigate);
            });
            item.appendChild(link);
        } else {
            item.textContent = segment;
        }
        
        breadcrumb.appendChild(item);
    });
}

// Find item by path
function findItemByPath(pathArray, root = fileSystem) {
    let current = root;
    
    for (const segment of pathArray) {
        if (!current.children) return null;
        const found = current.children.find(item => item.name === segment);
        if (!found) return null;
        current = found;
    }
    
    return current;
}

// Render file explorer
function renderFileExplorer() {
    fileExplorer.innerHTML = '';
    
    // Get current directory
    const currentDir = currentPath.length === 0 ? fileSystem : findItemByPath(currentPath);
    if (!currentDir || !currentDir.children) return;
    
    // Sort: directories first, then files, both alphabetically
    const sortedItems = [...currentDir.children].sort((a, b) => {
        if (a.type === b.type) {
            return a.name.localeCompare(b.name);
        }
        return a.type === 'directory' ? -1 : 1;
    });
    
    // Render each item
    sortedItems.forEach(item => {
        const itemElement = document.createElement('div');
        itemElement.className = `file-item ${item.type === 'directory' ? 'dir-item' : ''}`;
        
        const icon = document.createElement('i');
        icon.className = `bi bi-${getFileIcon(item.name, item.type)}`;
        
        const name = document.createElement('span');
        name.textContent = item.name;
        
        const size = document.createElement('span');
        size.className = 'file-size';
        size.textContent = item.type === 'file' ? formatFileSize(item.size || 0) : '';
        
        itemElement.appendChild(icon);
        itemElement.appendChild(name);
        itemElement.appendChild(size);
        
        itemElement.addEventListener('click', (e) => {
            e.stopPropagation();
            if (item.type === 'directory') {
                navigateTo([...currentPath, item.name]);
            } else {
                showFileContent(item);
            }
        });
        
        fileExplorer.appendChild(itemElement);
    });
    
    // Show empty state if no items
    if (sortedItems.length === 0) {
        const emptyState = document.createElement('div');
        emptyState.className = 'empty-state';
        emptyState.innerHTML = `
            <i class="bi bi-folder-x"></i>
            <p class="mb-0">This directory is empty</p>
        `;
        fileExplorer.appendChild(emptyState);
    }
}

// Show file content
function showFileContent(file) {
    fileTitle.textContent = file.name;
    fileActions.style.display = 'block';
    
    // Set up copy button
    copyBtn.onclick = () => {
        navigator.clipboard.writeText(file.content || '');
        const originalText = copyBtn.innerHTML;
        copyBtn.innerHTML = '<i class="bi bi-check"></i>';
        setTimeout(() => {
            copyBtn.innerHTML = originalText;
        }, 2000);
    };
    
    // Set up download button
    downloadBtn.onclick = () => {
        const blob = new Blob([file.content || ''], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = file.name;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    };
    
    // Determine content type and render accordingly
    let content = '';
    const ext = file.name.split('.').pop().toLowerCase();
    
    if (['jpg', 'jpeg', 'png', 'gif', 'svg'].includes(ext)) {
        // For images, create an image element
        content = `<div class="text-center"><img src="${file.path}" class="img-fluid" alt="${file.name}" onerror="this.parentElement.innerHTML='<i class=\'bi bi-file-earmark-image\'></i> <span>Cannot display image</span>';"></div>`;
    } else if (file.content) {
        // For text files, use highlight.js
        let codeClass = '';
        let codeContent = file.content;
        
        switch(ext) {
            case 'js':
                codeClass = 'javascript';
                break;
            case 'json':
                codeClass = 'json';
                try { codeContent = JSON.stringify(JSON.parse(file.content), null, 2); } catch(e) {}
                break;
            case 'html':
                codeClass = 'html';
                break;
            case 'css':
                codeClass = 'css';
                break;
            case 'md':
                // Render markdown
                content = `<div class="markdown-body">${window.marked.parse(file.content)}</div>`;
                break;
            default:
                codeClass = '';
        }
        
        if (content === '') {
            content = `<pre><code class="${codeClass}">${escapeHtml(codeContent)}</code></pre>`;
        }
    } else {
        content = '<div class="empty-state"><i class="bi bi-file-earmark"></i><p class="mb-0">No content to display</p></div>';
    }
    
    fileContent.innerHTML = content;
    
    // Apply syntax highlighting
    if (file.content && !['md', 'jpg', 'jpeg', 'png', 'gif', 'svg'].includes(ext)) {
        document.querySelectorAll('pre code').forEach((block) => {
            hljs.highlightElement(block);
        });
    }
    
    // Update active state
    document.querySelectorAll('.file-item').forEach(el => el.classList.remove('active'));
    const activeItem = Array.from(document.querySelectorAll('.file-item')).find(el => 
        el.textContent.trim().startsWith(file.name)
    );
    if (activeItem) activeItem.classList.add('active');
}

// Helper to escape HTML
function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Navigate to path
function navigateTo(pathArray) {
    currentPath = pathArray;
    renderBreadcrumb();
    renderFileExplorer();
    
    // Hide file content when navigating directories
    fileContent.innerHTML = '<div class="empty-state"><i class="bi bi-file-earmark-text"></i><p class="mb-0">Select a file to view its contents</p></div>';
    fileTitle.textContent = 'Select a file';
    fileActions.style.display = 'none';
    
    // Update URL hash
    const hash = currentPath.length > 0 ? '#' + currentPath.join('/') : '';
    window.history.pushState(null, '', window.location.pathname + hash);
}

// Handle back/forward navigation
window.addEventListener('popstate', () => {
    const hash = window.location.hash.slice(1);
    currentPath = hash ? hash.split('/') : [];
    renderBreadcrumb();
    renderFileExplorer();
});

// Initialize
async function init() {
    try {
        // Load file system from inline script or fallback to fetch
        if (typeof FILE_SYSTEM !== 'undefined') {
            fileSystem = FILE_SYSTEM;
        } else {
            const response = await fetch('data/file-system.json');
            if (!response.ok) throw new Error('Failed to load file system');
            fileSystem = await response.json();
        }
        
        // Handle initial navigation based on URL hash
        const hash = window.location.hash.slice(1);
        if (hash) {
            currentPath = hash.split('/').filter(Boolean);
            // Check if the path exists
            const item = findItemByPath(currentPath);
            if (!item) {
                currentPath = [];
                window.location.hash = '';
            } else if (item.type === 'file') {
                const dirPath = currentPath.slice(0, -1);
                currentPath = dirPath;
                renderBreadcrumb();
                renderFileExplorer();
                showFileContent(item);
                return;
            }
        }
        
        renderBreadcrumb();
        renderFileExplorer();
        
        // Set up refresh button
        refreshBtn.addEventListener('click', () => {
            window.location.reload();
        });
        
    } catch (error) {
        console.error('Error initializing file explorer:', error);
        fileExplorer.innerHTML = `
            <div class="alert alert-danger m-3">
                <i class="bi bi-exclamation-triangle-fill"></i>
                Failed to load file system: ${error.message}
            </div>
        `;
    }
}

// Start the application
document.addEventListener('DOMContentLoaded', init);

    const fileContent = document.getElementById('file-content');
    if (fileContent) {
        fileContent.innerHTML = `
            <div class="p-4 text-center text-danger">
                <i class="bi bi-exclamation-triangle-fill" style="font-size: 3rem;"></i>
                <h4 class="mt-3">Error</h4>
                <p class="mb-0">${message}</p>
            </div>`;
    }
}

// State
let currentPath = '';
let currentFile = null;

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Get file icon based on extension
function getFileIcon(filename) {
    const extension = filename.split('.').pop().toLowerCase();
    
    if (['jpg', 'jpeg', 'png', 'gif', 'svg'].includes(extension)) {
        return 'file-image';
    } else if (['js', 'jsx', 'ts', 'tsx', 'py', 'java', 'c', 'cpp', 'go', 'rs'].includes(extension)) {
        return 'file-code';
    } else if (['md', 'txt'].includes(extension)) {
        return 'file-text';
    } else if (['json', 'yaml', 'yml', 'xml'].includes(extension)) {
        return 'file-code';
    } else if (['pdf'].includes(extension)) {
        return 'file-pdf';
    } else if (['zip', 'tar', 'gz', 'rar'].includes(extension)) {
        return 'file-zip';
    } else {
        return 'file-earmark';
    }
}

// Render file content
function renderFileContent(content, filename) {
    const fileContent = document.getElementById('file-content');
    const extension = filename.split('.').pop().toLowerCase();
    
    if (extension === 'md') {
        fileContent.innerHTML = `
            <div class="p-4">
                <div class="markdown-body">
                    ${marked(content || 'No content available')}
                </div>
            </div>`;
    } else if (['js', 'jsx', 'ts', 'tsx', 'py', 'java', 'c', 'cpp', 'go', 'rs', 'css', 'html', 'json', 'yaml', 'yml', 'xml'].includes(extension)) {
        fileContent.innerHTML = `
            <div class="p-3">
                <pre><code class="language-${extension}">${content || ''}</code></pre>
            </div>`;
        hljs.highlightAll();
    } else {
        fileContent.innerHTML = `
            <div class="p-4">
                <pre class="mb-0">${content || 'No content available'}</pre>
            </div>`;
    }
}

// Update breadcrumb navigation
function updateBreadcrumb(path) {
    const breadcrumb = document.getElementById('breadcrumb');
    const parts = path ? path.split('/').filter(p => p) : [];
    
    let breadcrumbHtml = '<li class="breadcrumb-item"><a href="#" onclick="navigateTo(\'\')">Root</a></li>';
    
    let currentPath = '';
    parts.forEach((part, index) => {
        currentPath += (currentPath ? '/' : '') + part;
        const isLast = index === parts.length - 1;
        
        if (isLast) {
            breadcrumbHtml += `
                <li class="breadcrumb-item active" aria-current="page">${part}</li>`;
        } else {
            breadcrumbHtml += `
                <li class="breadcrumb-item"><a href="#" onclick="navigateTo('${currentPath}')">${part}</a></li>`;
        }
    });
    
    breadcrumb.innerHTML = breadcrumbHtml;
}

// Find node by path
function findNodeByPath(path) {
    if (!path) return FILE_SYSTEM;
    
    const parts = path.split('/').filter(p => p);
    let node = FILE_SYSTEM;
    
    for (const part of parts) {
        const found = node.children?.find(child => child.name === part);
        if (!found) return null;
        node = found;
    }
    
    return node;
}

// Render file explorer
function renderFileExplorer(path) {
    const explorer = document.getElementById('file-explorer');
    const node = findNodeByPath(path);
    
    if (!node || node.type !== 'directory') {
        explorer.innerHTML = `
            <div class="list-group-item text-center py-4 text-muted">
                <i class="bi bi-exclamation-triangle"></i>
                <div class="mt-2">Directory not found</div>
            </div>`;
        return;
    }
    
    let content = '';
    const dirs = [];
    const files = [];
    
    // Separate directories and files
    (node.children || []).forEach(child => {
        if (child.type === 'directory') {
            dirs.push(child);
        } else {
            files.push(child);
        }
    });
    
    // Sort directories and files
    dirs.sort((a, b) => a.name.localeCompare(b.name));
    files.sort((a, b) => a.name.localeCompare(b.name));
    
    // Add parent directory link if not in root
    if (path) {
        const parentPath = path.split('/').slice(0, -1).join('/');
        content += `
            <div class="list-group-item list-group-item-action" onclick="navigateTo('${parentPath}')">
                <i class="bi bi-folder2"></i> ..
            </div>`;
    }
    
    // Add directories
    dirs.forEach(dir => {
        content += `
            <div class="list-group-item list-group-item-action dir-item" 
                 onclick="navigateTo('${dir.path}')">
                <i class="bi bi-folder2"></i> ${dir.name}
            </div>`;
    });
    
    // Add files
    files.forEach(file => {
        const icon = getFileIcon(file.name);
        const isActive = currentFile === file.path ? 'active' : '';
        content += `
            <div class="list-group-item list-group-item-action ${isActive}" 
                 onclick="openFile('${file.path}')" 
                 data-path="${file.path}">
                <i class="bi bi-${icon}"></i> ${file.name}
                <span class="float-end text-muted small">${formatFileSize(file.size)}</span>
            </div>`;
    });
    
    explorer.innerHTML = content || `
        <div class="list-group-item text-center py-4 text-muted">
            <i class="bi bi-folder-x"></i>
            <div class="mt-2">Directory is empty</div>
        </div>`;
    
    // Update file count
    const fileCount = document.getElementById('file-count');
    if (dirs.length === 0 && files.length === 0) {
        fileCount.textContent = 'Empty';
    } else {
        fileCount.textContent = `${dirs.length} ${dirs.length === 1 ? 'folder' : 'folders'}, ${files.length} ${files.length === 1 ? 'file' : 'files'}`;
    }
}

// Navigate to directory
function navigateTo(path) {
    currentPath = path;
    currentFile = null;
    renderFileExplorer(path);
    updateBreadcrumb(path);
    
    // Update file content view
    const fileContent = document.getElementById('file-content');
    fileContent.innerHTML = `
        <div class="p-4 text-center text-muted">
            <i class="bi bi-folder2-open" style="font-size: 3rem; opacity: 0.5;"></i>
            <h4 class="mt-3">${path ? path.split('/').pop() || 'Root' : 'Root'}</h4>
            <p class="mb-0">Select a file to view its contents</p>
        </div>`;
}

// Open file
function openFile(path) {
    const node = findNodeByPath(path);
    if (!node || node.type !== 'file') return;
    
    currentFile = path;
    
    // Update active state in file list
    document.querySelectorAll('.file-item').forEach(item => {
        item.classList.remove('active');
        if (item.getAttribute('data-path') === path) {
            item.classList.add('active');
        }
    });
    
    // Render file content
    renderFileContent(node.content, node.name);
}

// Initialize file explorer
async function initFileExplorer() {
    // Show loading state
    const explorer = document.getElementById('file-explorer');
    if (explorer) {
        explorer.innerHTML = `
            <div class="list-group-item text-center py-4">
                <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
                <div class="mt-2 text-muted">Loading file system...</div>
            </div>`;
    }

    // Load file system
    await loadFileSystem();
    
    if (!FILE_SYSTEM) {
        return; // Error already shown by loadFileSystem
    }
    
    // Make functions available globally
    window.navigateTo = navigateTo;
    window.openFile = openFile;
    
    // Set up download button
    document.getElementById('download-btn').addEventListener('click', function() {
        if (currentFile) {
            const node = findNodeByPath(currentFile);
            if (node && node.type === 'file') {
                const blob = new Blob([node.content || ''], { type: 'text/plain' });
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = node.name;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
            }
        }
    });
    
    // Set up copy button
    document.getElementById('copy-btn').addEventListener('click', function() {
        if (currentFile) {
            const node = findNodeByPath(currentFile);
            if (node && node.type === 'file') {
                navigator.clipboard.writeText(node.content || '');
                const btn = this;
                const originalHTML = btn.innerHTML;
                btn.innerHTML = '<i class="bi bi-check"></i>';
                setTimeout(() => {
                    btn.innerHTML = originalHTML;
                }, 2000);
            }
        }
    });
    
    // Initial render
    navigateTo('');
}

// Start the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', initFileExplorer);
