// The entry file of your WebAssembly module.
import { Context, generateEvent, Storage } from '@massalabs/massa-as-sdk';
import { Args } from '@massalabs/as-types';

/**
 * This function is meant to be called only one time: when the contract is deployed.
 *
 * @param binaryArgs - Arguments serialized with Args
 */

const COUNTER_KEY = "counter";

export function constructor() {
  // This line is important. It ensures that this function can't be called in the future.
  // If you remove this check, someone could call your constructor function and reset your smart contract.
  if (!Context.isDeployingContract()) {
    return;
  }

  Storage.set(COUNTER_KEY, '0');
  
  generateEvent(`Constructor Counter Smart Contract`);
}

export function incrementByOne(): u32 {
    const currentValue = readCounter();
    Storage.set(COUNTER_KEY, (currentValue + 1).toString());
    return readCounter();
}

export function incrementByN(args: Args): u32 {
    const n = args.nextU32().unwrap();
    const currentValue = readCounter();
    Storage.set(COUNTER_KEY, (currentValue + n).toString());
    return readCounter();
}

export function readCounter(): u32 {
    const value = Storage.get(COUNTER_KEY);
    return value ? parseInt(value) : 0;
}