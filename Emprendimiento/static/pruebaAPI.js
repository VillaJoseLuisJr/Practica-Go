const btn = document.getElementById('btnFetchCharacters');
const div = document.getElementById('prueba-API')

btn.addEventListener('click', () => {
    console.log('Fetch API');

    fetch('https://rickandmortyapi.com/api/character')
        .then((response) => response.json())
        .then((data) => renderCharacters(data.results));

});

function renderCharacters(characters) {
    const contenedor = document.getElementById('prueba-API');
    contenedor.innerHTML = '';

    characters.forEach(ch => {
        const div = document.createElement('div');
        div.className = 'producto';

        // Crear la sección de texto
        const divTexto = document.createElement('div');
        divTexto.className = 'div-texto';
        divTexto.innerHTML = `
                    <h2>Nombre: ${ch.name}</h2>
                    <p>Status: ${ch.status}</p>
                    <p>El ID de este personaje es: ${ch.id}</p>
                `;

        // Crear la sección de imagen
        const divImagenes = document.createElement('div');
        divImagenes.className = 'div-imagenes';
        const img = document.createElement('img');
        img.src = ch.image;
        divImagenes.appendChild(img);



        // Ensamblar todo
        div.appendChild(divTexto)
        div.appendChild(divImagenes)
        contenedor.appendChild(div);
    });
}