function showMessage(message) {
    $.toaster({
        message: message ? message : 'Данные сохранены', title: '',
        settings: {
            'timeout': 2000,
            'toaster': {
                'css':
                    {
                        'position': 'fixed',
                        'top': '50px',
                        'left': '50%',
                        'width': '300px',
                        'zIndex': 50000
                    }
            },
            'toast':
                {
                    'template':
                        '<div class="alert alert-%priority% alert-dismissible" role="alert">' +

                        '<span class="message"></span>' +
                        '</div>',
                    'css': {},
                    'cssm': {},
                    'csst': {'fontWeight': 'bold'},
                    'fade': 'slow',


                },
        },


    });
}

function showError(message, title, top, left) {
    $.toaster({
        message: message, title: title ? title : 'Ошибка', priority: 'danger',
        settings: {
            'timeout': 3000,
            'toaster': {
                'css':
                    {
                        'position': 'fixed',
                        'top': top ? top : '50px',
                        'left': left ? left : '50%',
                        'width': '300px',
                        'zIndex': 50000
                    }
            },
            'toast':
                {
                    'template':
                        '<div class="alert alert-%priority% alert-dismissible" role="alert">' +
                        '<span class="title"> </span><br/>' +
                        '<span class="message"></span>' +
                        '</div>',
                    'css': {},
                    'cssm': {},
                    'csst': {'fontWeight': 'bold'},
                    'fade': 'slow',


                },
        },


    });
}