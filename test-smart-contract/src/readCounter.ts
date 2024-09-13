import * as dotenv from 'dotenv';
import { getEnvVariable } from './utils';
import { Client, bytesToStr, ClientFactory, DefaultProviderUrls } from '@massalabs/massa-web3';
import { WalletClient } from '@massalabs/massa-sc-deployer';
import {
  Args,
  fromMAS,
  MAX_GAS_DEPLOYMENT,
  CHAIN_ID,
} from '@massalabs/massa-web3';

// Load .env file content into process.env
dotenv.config();

// Get environment variables
const publicApi = getEnvVariable('JSON_RPC_URL_PUBLIC');
const secretKey = getEnvVariable('WALLET_SECRET_KEY');
const contractAddress = getEnvVariable('CONTRACT_ADDRESS')
// Define deployment parameters
const chainId = CHAIN_ID.BuildNet; // Choose the chain ID corresponding to the network you want to deploy to
const maxGas = MAX_GAS_DEPLOYMENT; // Gas for deployment Default is the maximum gas allowed for deployment
const fees = 10000000n; // Fees to be paid for deployment. Default is 0
const waitFirstEvent = true;

// Create an account using the private keyc
const baseAccount = await WalletClient.getAccountFromSecretKey(secretKey);

// const client = new Client({providers: new IPro`publicApi});

async function initMassaBuildnetClient() {
    try {  
      // Create a client using the buildnet nodes
      const client: Client = await ClientFactory.createDefaultClient(
        DefaultProviderUrls.BUILDNET,
        CHAIN_ID.BuildNet,
        true, // retry failed requests
        baseAccount // optional parameter
      );

      console.log('Massa Web3 buildnet client initialized successfully');
      return client;
    } catch (error) {
      console.error('Failed to initialize Massa Web3 buildnet client:', error);
    }
}

async function main() {
    const client = await initMassaBuildnetClient();

    const res = await client?.smartContracts().readSmartContract({
        maxGas: BigInt(2100000),
        targetAddress: contractAddress,
        targetFunction: "readCounter",
        parameter: [],
    });
      
    const count = res?.returnValue? bytesToStr(res?.returnValue): "";

    console.log(count)
}

(async () => {
  await main();
  process.exit(0); // terminate the process after deployment(s)
})();
