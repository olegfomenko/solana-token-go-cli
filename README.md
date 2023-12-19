# Solana Token CLI

CLI that helps to create your token.

Go through the following steps:

1. Install Solana CLI: "<https://docs.solana.com/cli/install-solana-cli-tools>"

2. Create filesystem wallet: "<https://docs.solana.com/wallet-guide/file-system-wallet>"

3. Fund it with some SOL tokens to pay commissions

4. Read the "<https://spl.solana.com/token>" article. Below, the commands to execute with short description will be provided.

5. Check your keypair file path: `solana config get`

6. Install cli: `cargo install spl-token-cli`

7. Create token: `spl-token create-token --decimals 9`. For NFT set decimals to 0. Save the token address.

8. Create account to hold tokens: `spl-token create-account <token account>`

9. Mint some tokens `spl-token mint <token account> 100 <created holder account>`. Mint only 1 if it should be an NFT.

10. Provide environment var: `SENDER_PRV=/path/to/keyfile`. Also, you can setup `SOLANA_RPC` env with custom RPC address (if empty - default will be used).

11. Setup metadata file, example: [here](./metadata.json).

12. Use Docker image or binary file from the latest release to execute the following command: `./solana-token-cli [devnet|testnet|mainnet] [path to metadata file] [token address]`



