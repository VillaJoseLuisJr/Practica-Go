async function fetchData() {
    try {
        const nombrePokemon = document.getElementById("nombrePokemon").value.toLowerCase();

        const response = await fetch(`https://pokeapi.co/api/v2/pokemon/${nombrePokemon}`);

        if (!response.ok) {
            throw new Error("No se pudo hacer fetch al recurso");
        }
        const data = await response.json();
        const pokemonSprite = data.sprites.front_default;
        const imgElement = document.getElementById("pokemonSprite");

        imgElement.src = pokemonSprite;
        imgElement.style.display = "block";
    }
    catch (error) {
        console.log(error)
    }
}