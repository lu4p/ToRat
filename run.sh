#!/bin/bash
docker run -it -v "$(pwd)"/dist:/dist_ext -p 8000:8000 torat

