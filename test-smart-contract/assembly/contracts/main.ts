// The entry file of your WebAssembly module.
import { Context, generateEvent, Storage } from '@massalabs/massa-as-sdk';
import { Args } from '@massalabs/as-types';

const COUNTER_KEY = "counter";

export function constructor(): void {
  // This line is important. It ensures that this function can't be called in the future.
  // If you remove this check, someone could call your constructor function and reset your smart contract.
  if (!Context.isDeployingContract()) {
    return;
  }

  Storage.set(COUNTER_KEY, "0");
  
  generateEvent(`Constructor Counter Smart Contract`);
}

export function incrementByOne(): void {
  const currentValue = parseInt(Storage.get(COUNTER_KEY));
  const value = currentValue + 1;
  Storage.set(COUNTER_KEY, value.toString());
}

export function incrementByN(binaryArgs: StaticArray<u8>): void {
  const args = new Args(binaryArgs)
  const currentValue = parseInt(Storage.get(COUNTER_KEY));
  const value = currentValue + args.nextU32().unwrap();
  Storage.set(COUNTER_KEY, value.toString());
}

export function readCounter(): string {
  const value = Storage.get(COUNTER_KEY);
  return value;
}