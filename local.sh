alias swapd=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .swap
        shift
        ;;
    esac
done

if ! [ -f .swap/data/priv_validator_state.json ]; then
  swapd init validator --chain-id "swap-1" --home .swap &> /dev/null

  swapd keys add validator --home .swap --keyring-backend test &> /dev/null
  swapd genesis add-genesis-account validator 1000000ustake --home .swap --keyring-backend test
  swapd keys add authority --recover --home .swap --keyring-backend test <<< "occur subway woman achieve deputy rapid museum point usual appear oil blue rate title claw debate flag gallery level object baby winner erase carbon" &> /dev/null
  swapd genesis add-genesis-account authority 10000000uusdc --home .swap --keyring-backend test
  swapd keys add provider --home .swap --keyring-backend test &> /dev/null
  swapd genesis add-genesis-account provider 10000000uusdn --home .swap --keyring-backend test

  swapd keys add user --recover --home .swap --keyring-backend test <<< "entry cake clinic beach able model all doll combine kit sausage essay" &> /dev/null
  swapd genesis add-genesis-account user 10000000000000uusdn,10000000000000uusdc,10000000000000uusde --home .swap --keyring-backend test

  TEMP=.swap/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .swap/config/genesis.json > $TEMP && mv $TEMP .swap/config/genesis.json

  swapd genesis gentx validator 1000000ustake --chain-id "swap-1" --home .swap --keyring-backend test &> /dev/null
  swapd genesis collect-gentxs --home .swap &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .swap/config/config.toml
fi

swapd start --home .swap
