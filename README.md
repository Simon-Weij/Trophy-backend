# Trophy

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

Trophy is an open source and self-hostable clipping software. This repository belongs to the self-hostable backend. For the desktop-app see [this](https://github.com/Simon-Weij/Trophy). want to help? See [contributing](#contributing)

## Contributing

Contributions are always welcome! If you plan on implementing a feature, please open an issue beforehand. Same with large bugs. For small bugs like typoes, just open a PR.

## Deployment

The easiest way to deploy Trophy is using [docker compose](https://docs.docker.com/compose/), podman may work, but has not been tested.

Currently there are no pre-build images, so building yourself is the only supported way.

```bash
# Clone the repo
git clone https://github.com/Simon-Weij/Trophy

# cd into the directory
cd Trophy

# Copy .env.example to .env
cp .env.example .env

# Fill out .env with your secure credentials
nano .env

# Run it!
docker compose up
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
