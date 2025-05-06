document.getElementById('buscador').addEventListener('submit', function (e) {
    e.preventDefault();

    const query = document.getElementById('busqueda').value;

    fetch('/buscar?busqueda=' + encodeURIComponent(query))
        .then(response => response.json())
        .then(productos => {
            const contenedor = document.getElementById('productos');
            contenedor.innerHTML = '';

            if (productos.length === 0) {
                contenedor.innerHTML = '<p>No se encontraron productos.<p>';
                return;
            }

            console.log(productos)
            productos.forEach(p => {
                const div = document.createElement('div');
                div.className = 'producto';

                // Crear la sección de texto
                const divTexto = document.createElement('div');
                divTexto.className = 'div-texto';
                divTexto.innerHTML = `
                    <h2>${p.nombre}</h2>
                    <p>${p.descripcion}</p>
                    <p>El ID de este producto es: ${p.id}</p>
                `;

                // Crear la sección de imagen
                const divImagenes = document.createElement('div');
                divImagenes.className = 'div-imagenes';
                if (p.imagen) {
                    const img = document.createElement('img');
                    img.src = p.imagen;
                    img.alt = p.nombre;
                    divImagenes.appendChild(img);
                } else {
                    divImagenes.innerHTML = `
                        <svg xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 512 512"><!--!Font Awesome Free 6.7.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2025 Fonticons, Inc.-->
                        <path
                            d="M0 96C0 60.7 28.7 32 64 32l384 0c35.3 0 64 28.7 64 64l0 320c0 35.3-28.7 64-64 64L64 480c-35.3 0-64-28.7-64-64L0 96zM323.8 202.5c-4.5-6.6-11.9-10.5-19.8-10.5s-15.4 3.9-19.8 10.5l-87 127.6L170.7 297c-4.6-5.7-11.5-9-18.7-9s-14.2 3.3-18.7 9l-64 80c-5.8 7.2-6.9 17.1-2.9 25.4s12.4 13.6 21.6 13.6l96 0 32 0 208 0c8.9 0 17.1-4.9 21.2-12.8s3.6-17.4-1.4-24.7l-120-176zM112 192a48 48 0 1 0 0-96 48 48 0 1 0 0 96z" />
                        </svg>
                        <p><em>No hay imagen</em></p>
                    `;
                }

                // Crear el botón de eliminar
                const divBoton = document.createElement('div');
                divBoton.className = 'div-boton';
                divBoton.innerHTML = `
                    <form action="/eliminar" method="post"
                        onsubmit="return confirm('¿Estás seguro de que querés eliminar este producto?');">
                        <input type="hidden" name="id" value="${p.id}">
                        <button type="submit" class="boton">Eliminar <svg xmlns="http://www.w3.org/2000/svg"
                                id="tacho-basura"
                                viewBox="0 0 448 512"><!--!Font Awesome Free 6.7.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2025 Fonticons, Inc.-->
                                <path
                                    d="M135.2 17.7C140.6 6.8 151.7 0 163.8 0L284.2 0c12.1 0 23.2 6.8 28.6 17.7L320 32l96 0c17.7 0 32 14.3 32 32s-14.3 32-32 32L32 96C14.3 96 0 81.7 0 64S14.3 32 32 32l96 0 7.2-14.3zM32 128l384 0 0 320c0 35.3-28.7 64-64 64L96 512c-35.3 0-64-28.7-64-64l0-320zm96 64c-8.8 0-16 7.2-16 16l0 224c0 8.8 7.2 16 16 16s16-7.2 16-16l0-224c0-8.8-7.2-16-16-16zm96 0c-8.8 0-16 7.2-16 16l0 224c0 8.8 7.2 16 16 16s16-7.2 16-16l0-224c0-8.8-7.2-16-16-16zm96 0c-8.8 0-16 7.2-16 16l0 224c0 8.8 7.2 16 16 16s16-7.2 16-16l0-224c0-8.8-7.2-16-16-16z" />
                            </svg></button>
                    </form>
                `;

                // Ensamblar todo
                div.appendChild(divTexto)
                div.appendChild(divImagenes)
                div.appendChild(divBoton)
                contenedor.appendChild(div);
            });
        })
        .catch(error => {
            console.error('Error en la búsqueda:', error);
        });
});