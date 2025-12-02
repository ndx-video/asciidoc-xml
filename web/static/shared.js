// ./web/static/shared.js
// Shared navbar and footer components

// Global error handler - must be defined first
(function() {
    'use strict';
    
    // Function to send errors to the server
    function reportErrorToServer(errorInfo) {
        try {
            fetch('/api/jserror', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(errorInfo)
            }).catch(function(err) {
                // Silent failure - don't create error loops
                console.error('Failed to report error to server:', err);
            });
        } catch (e) {
            // Silent failure
            console.error('Exception in reportErrorToServer:', e);
        }
    }
    
    // Global error handler
    window.addEventListener('error', function(event) {
        try {
            var errorInfo = {
                message: event.message || 'Unknown error',
                stack: event.error ? event.error.stack : '',
                url: event.filename || window.location.href,
                lineNumber: event.lineno || 0,
                colNumber: event.colno || 0,
                userAgent: navigator.userAgent,
                timestamp: new Date().toISOString(),
                location: window.location.href
            };
            
            console.error('JavaScript Error:', errorInfo);
            reportErrorToServer(errorInfo);
        } catch (e) {
            console.error('Error in error handler:', e);
        }
    });
    
    // Global unhandled promise rejection handler
    window.addEventListener('unhandledrejection', function(event) {
        try {
            var errorInfo = {
                message: 'Unhandled Promise Rejection: ' + (event.reason ? event.reason.message || event.reason : 'Unknown'),
                stack: event.reason && event.reason.stack ? event.reason.stack : '',
                url: window.location.href,
                lineNumber: 0,
                colNumber: 0,
                userAgent: navigator.userAgent,
                timestamp: new Date().toISOString(),
                location: window.location.href
            };
            
            console.error('Unhandled Promise Rejection:', errorInfo);
            reportErrorToServer(errorInfo);
        } catch (e) {
            console.error('Error in unhandledrejection handler:', e);
        }
    });
    
    // Make reportErrorToServer available globally for manual error reporting
    window.reportErrorToServer = reportErrorToServer;
})();

(function() {
    'use strict';
    
    function initSharedComponents() {
        try {
            // Footer HTML
            const footerHTML = `
            <footer>
                PicoCSS embedded in HTML output • goja used for testing only • This project is brought to you by NDX Pty Ltd. Contributions are welcome.
            </footer>
        `;
        
        // Process navbar - enhance existing structure
        const navbarPlaceholder = document.getElementById('navbar-placeholder');
        if (navbarPlaceholder) {
            const existingNav = navbarPlaceholder.querySelector('.navbar');
            if (existingNav) {
                // Ensure navbar-content wrapper exists
                let navbarContent = existingNav.querySelector('.navbar-content');
                if (!navbarContent) {
                    const content = existingNav.innerHTML;
                    existingNav.innerHTML = `<div class="navbar-content">${content}</div>`;
                    navbarContent = existingNav.querySelector('.navbar-content');
                }
                
                // Move non-menu items out of navbar (output-selector, buttons, status)
                // These should be in a control strip (home page) or removed from navbar (other pages)
                const outputSelector = navbarContent.querySelector('.output-selector');
                const buttons = Array.from(navbarContent.querySelectorAll('button:not(.navbar-toggle)'));
                const status = navbarContent.querySelector('.status');
                
                // Check if we're on the home page (has output-selector or multiple buttons)
                const isHomePage = outputSelector !== null || buttons.length > 2;
                
                // Remove these from navbar
                if (outputSelector) {
                    outputSelector.remove();
                }
                buttons.forEach(btn => btn.remove());
                if (status) {
                    status.remove();
                }
                
                // Add burger menu if not present - insert after title, before menu
                if (!navbarContent.querySelector('.navbar-toggle')) {
                    const toggle = document.createElement('button');
                    toggle.className = 'navbar-toggle';
                    toggle.id = 'navbar-toggle';
                    toggle.setAttribute('aria-label', 'Toggle menu');
                    toggle.innerHTML = '<span></span><span></span><span></span>';
                    const menu = navbarContent.querySelector('.navbar-menu');
                    if (menu) {
                        navbarContent.insertBefore(toggle, menu);
                    } else {
                        const title = navbarContent.querySelector('.navbar-title, h1');
                        if (title) {
                            navbarContent.insertBefore(toggle, title.nextSibling);
                        } else {
                            navbarContent.insertBefore(toggle, navbarContent.firstChild);
                        }
                    }
                }
                
                // Ensure menu items are wrapped in navbar-menu
                const menu = navbarContent.querySelector('.navbar-menu');
                if (!menu) {
                    const links = Array.from(navbarContent.querySelectorAll('a.navbar-link'));
                    if (links.length > 0) {
                        const menuDiv = document.createElement('div');
                        menuDiv.className = 'navbar-menu';
                        menuDiv.id = 'navbar-menu';
                        links.forEach(link => {
                            if (link.parentNode === navbarContent) {
                                menuDiv.appendChild(link);
                            }
                        });
                        const toggle = navbarContent.querySelector('.navbar-toggle');
                        if (toggle) {
                            navbarContent.insertBefore(menuDiv, toggle.nextSibling);
                        } else {
                            const title = navbarContent.querySelector('.navbar-title, h1');
                            if (title) {
                                navbarContent.insertBefore(menuDiv, title.nextSibling);
                            }
                        }
                    }
                }
                
                // Create control strip for home page if needed
                if (isHomePage && (outputSelector || buttons.length > 0 || status)) {
                    let controlStrip = document.getElementById('home-control-strip');
                    if (!controlStrip) {
                        controlStrip = document.createElement('div');
                        controlStrip.id = 'home-control-strip';
                        controlStrip.className = 'home-control-strip';
                        const navbar = document.querySelector('.navbar');
                        if (navbar && navbar.nextSibling) {
                            navbar.parentNode.insertBefore(controlStrip, navbar.nextSibling);
                        } else if (navbar) {
                            navbar.parentNode.appendChild(controlStrip);
                        }
                    }
                    
                    // Clear existing content
                    controlStrip.innerHTML = '';
                    
                    // Add elements to control strip
                    if (outputSelector) {
                        controlStrip.appendChild(outputSelector);
                    }
                    buttons.forEach(btn => controlStrip.appendChild(btn));
                    if (status) {
                        controlStrip.appendChild(status);
                    }
                    
                    // Show the control strip
                    controlStrip.style.display = 'flex';
                } else if (!isHomePage && status) {
                    // For other pages, move status to a less prominent location or remove from navbar
                    // Status will be removed from navbar but not added to control strip
                }
            }
        }
        
        // Inject footer
        const footerPlaceholder = document.getElementById('footer-placeholder');
        if (footerPlaceholder) {
            // Replace the placeholder with the actual footer
            footerPlaceholder.outerHTML = footerHTML;
        } else {
            // Try to find existing footer and replace
            const existingFooter = document.querySelector('footer');
            if (existingFooter && existingFooter.id !== 'footer-placeholder') {
                existingFooter.outerHTML = footerHTML;
            } else {
                // Create footer at end of body if nothing exists
                const body = document.body;
                if (body) {
                    body.insertAdjacentHTML('beforeend', footerHTML);
                }
            }
        }
        
        // Burger menu toggle
        function initBurgerMenu() {
            try {
                const toggle = document.getElementById('navbar-toggle');
                const menu = document.getElementById('navbar-menu');
                if (toggle && menu) {
                    toggle.addEventListener('click', (e) => {
                        try {
                            e.stopPropagation();
                            menu.classList.toggle('active');
                            toggle.classList.toggle('active');
                        } catch (error) {
                            console.error('Error in burger menu click handler:', error);
                        }
                    });
                    
                    // Close menu when clicking outside
                    document.addEventListener('click', (e) => {
                        try {
                            if (!toggle.contains(e.target) && !menu.contains(e.target)) {
                                menu.classList.remove('active');
                                toggle.classList.remove('active');
                            }
                        } catch (error) {
                            console.error('Error in document click handler:', error);
                        }
                    });
                }
            } catch (error) {
                console.error('Error in initBurgerMenu:', error);
            }
        }
        
            initBurgerMenu();
        } catch (error) {
            console.error('Error in initSharedComponents:', error);
            if (window.reportErrorToServer) {
                window.reportErrorToServer({
                    message: 'Error in initSharedComponents: ' + error.message,
                    stack: error.stack || '',
                    url: window.location.href,
                    lineNumber: 0,
                    colNumber: 0,
                    userAgent: navigator.userAgent,
                    timestamp: new Date().toISOString(),
                    location: window.location.href
                });
            }
        }
    }
    
    // Initialize when DOM is ready
    try {
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', function() {
                try {
                    initSharedComponents();
                } catch (error) {
                    console.error('Error during DOMContentLoaded:', error);
                }
            });
        } else {
            initSharedComponents();
        }
    } catch (error) {
        console.error('Error setting up initialization:', error);
    }
})();
