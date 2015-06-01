#!/bin/sh

set -x
set -e



# AWE server (miminal image)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server/shock -XPUT -d value="shock.metagenomics.anl.gov/node/a074d9fe-5c8e-4424-987b-d8ffc96da618"

# AWE client (miminal image)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-client/shock -XPUT -d value="shock.metagenomics.anl.gov/node/0ed11256-0067-474a-b365-d5c951433211"

# MongoDB (for AWE server)
curl -L http://127.0.0.1:4001/v2/keys/service_images/awe-server-mongodb/shock -XPUT -d value="shock.metagenomics.anl.gov/node/6dbd1649-0ad2-4c44-887b-aafeb02849fa"
