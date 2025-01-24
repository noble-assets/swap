find types -type f -name "*.pb.go" -exec rm -f {} \;
find types -type f -name "*.pb.gw.go" -exec rm -f {} \;

cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..


cp -r swap.noble.xyz/* ./
rm -rf swap.noble.xyz
cp -r api/noble/swap/* api/
find api/ -type f -name "*.go" -exec sed -i 's|swap.noble.xyz/api/noble/swap|swap.noble.xyz/api|g' {} +

rm -rf api/noble
rm -rf noble
