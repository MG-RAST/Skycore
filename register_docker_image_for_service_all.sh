#!/bin/sh

set -x
set -e

# solr-m5nr
curl -L http://127.0.0.1:4001/v2/keys/service_images/solr-m5nr/shock -XPUT -d value="shock.metagenomics.anl.gov/node/174bbd39-c80c-4473-964b-2b97a226d10c"

# mg-rast-v4-web
curl -L http://127.0.0.1:4001/v2/keys/service_images/mg-rast-v4-web/shock -XPUT -d value="shock.metagenomics.anl.gov/node/247d49e8-5699-4329-92cc-774a210b8dff"

# AWE server
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server/shock -XPUT -d value="shock.metagenomics.anl.gov/node/a074d9fe-5c8e-4424-987b-d8ffc96da618"

# AWE client
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-client/shock -XPUT -d value="shock.metagenomics.anl.gov/node/4cdfad56-ed3c-40d3-87cc-b97d4fb7588d"

# MongoDB (for AWE server)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server-mongodb/shock -XPUT -d value="shock.metagenomics.anl.gov/node/6dbd1649-0ad2-4c44-887b-aafeb02849fa"