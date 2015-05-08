#!/bin/sh

set -x
set -e


# mg-rast-nginx (note that nginx and confd use the same image)
curl -L http://127.0.0.1:4001/v2/keys/service_images/mg-rast-nginx/shock -XPUT -d value="shock.metagenomics.anl.gov/node/acdd0450-22e2-4ead-8956-dee1cbfbc8b9"

# mg-rast-confd (note that nginx and confd use the same image)
curl -L http://127.0.0.1:4001/v2/keys/service_images/mg-rast-confd/shock -XPUT -d value="shock.metagenomics.anl.gov/node/acdd0450-22e2-4ead-8956-dee1cbfbc8b9"


# solr-m5nr
curl -L http://127.0.0.1:4001/v2/keys/service_images/solr-m5nr/shock -XPUT -d value="shock.metagenomics.anl.gov/node/82129683-8689-43bb-a9ec-4d4bdeb6248a"

# mg-rast-v4-web
curl -L http://127.0.0.1:4001/v2/keys/service_images/mg-rast-v4-web/shock -XPUT -d value="shock.metagenomics.anl.gov/node/247d49e8-5699-4329-92cc-774a210b8dff"

# AWE server (miminal image)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server/shock -XPUT -d value="shock.metagenomics.anl.gov/node/a074d9fe-5c8e-4424-987b-d8ffc96da618"

# AWE client (miminal image)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-client/shock -XPUT -d value="shock.metagenomics.anl.gov/node/0ed11256-0067-474a-b365-d5c951433211"

# MongoDB (for AWE server)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server-mongodb/shock -XPUT -d value="shock.metagenomics.anl.gov/node/6dbd1649-0ad2-4c44-887b-aafeb02849fa"
