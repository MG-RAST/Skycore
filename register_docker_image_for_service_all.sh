#!/bin/sh

set -x
set -e

# solr-m5nr
curl -L http://127.0.0.1:4001/v2/keys/service_images/solr-m5nr/shock -XPUT -d value="shock.metagenomics.anl.gov/node/174bbd39-c80c-4473-964b-2b97a226d10c"

# mg-rast-v4-web
curl -L http://127.0.0.1:4001/v2/keys/service_images/mg-rast-v4-web/shock -XPUT -d value="shock.metagenomics.anl.gov/node/247d49e8-5699-4329-92cc-774a210b8dff"

# AWE server
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server/shock -XPUT -d value="shock.metagenomics.anl.gov/node/e88a13a6-ab36-4ae3-aab2-5f09476f9556"

# AWE client
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-client/shock -XPUT -d value="shock.metagenomics.anl.gov/node/49186e36-6e4c-4606-b4ca-86c9ebf28c49"

# MongoDB (for AWE server)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server-mongodb/shock -XPUT -d value="shock.metagenomics.anl.gov/node/6dbd1649-0ad2-4c44-887b-aafeb02849fa"