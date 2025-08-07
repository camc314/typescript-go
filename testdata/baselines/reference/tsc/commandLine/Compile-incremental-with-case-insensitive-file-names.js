currentDirectory::/home/project
useCaseSensitiveFileNames::false
Input::
//// [/home/node_modules/lib1/index.d.ts] *new* 
import type { Foo } from 'someLib';
export type { Foo as Foo1 };
//// [/home/node_modules/lib1/package.json] *new* 
{
    "name": "lib1"
}
//// [/home/node_modules/lib2/index.d.ts] *new* 
import type { Foo } from 'somelib';
export type { Foo as Foo2 };
export declare const foo2: Foo;
//// [/home/node_modules/lib2/package.json] *new* 
{
    "name": "lib2"
}
//// [/home/node_modules/otherLib/index.d.ts] *new* 
export type Str = string;
//// [/home/node_modules/otherLib/package.json] *new* 
{
    "name": "otherlib"
}
//// [/home/node_modules/someLib/index.d.ts] *new* 
import type { Str } from 'otherLib';
export type Foo = { foo: Str; };
//// [/home/node_modules/someLib/package.json] *new* 
        {
"name": "somelib"
        }
//// [/home/project/src/index.ts] *new* 
import type { Foo1 } from 'lib1';
import type { Foo2 } from 'lib2';
export const foo1: Foo1 = { foo: "a" };
export const foo2: Foo2 = { foo: "b" };
//// [/home/project/tsconfig.json] *new* 
        {
"compilerOptions": {
				"incremental": true
},
        }

tsgo -p .
ExitStatus:: Success
Output::
//// [/home/project/src/index.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo2 = exports.foo1 = void 0;
exports.foo1 = { foo: "a" };
exports.foo2 = { foo: "b" };

//// [/home/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../src/tslibs/ts/lib/lib.d.ts","../node_modules/otherlib/index.d.ts","../node_modules/somelib/index.d.ts","../node_modules/lib1/index.d.ts","../node_modules/lib2/index.d.ts","./src/index.ts"],"fileInfos":[{"version":"8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };","affectsGlobalScope":true,"impliedNodeFormat":1},"1fe659ed0634bb57b6dc25e9062f1162-export type Str = string;","12e112ff6e2744bb42d8e0b511e44117-import type { Str } from 'otherLib';\nexport type Foo = { foo: Str; };","b6305455d920a6729c435e6acf45eff6-import type { Foo } from 'someLib';\nexport type { Foo as Foo1 };","a5393e550a9c20a242a120bf6410db48-import type { Foo } from 'somelib';\nexport type { Foo as Foo2 };\nexport declare const foo2: Foo;","42aef197ff5f079223e2c29fb2e77cc5-import type { Foo1 } from 'lib1';\nimport type { Foo2 } from 'lib2';\nexport const foo1: Foo1 = { foo: \"a\" };\nexport const foo2: Foo2 = { foo: \"b\" };"],"fileIdsList":[[3],[2],[4,5]],"referencedMap":[[4,1],[5,1],[3,2],[6,3]]}
//// [/home/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../src/tslibs/ts/lib/lib.d.ts",
    "../node_modules/otherlib/index.d.ts",
    "../node_modules/somelib/index.d.ts",
    "../node_modules/lib1/index.d.ts",
    "../node_modules/lib2/index.d.ts",
    "./src/index.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../src/tslibs/ts/lib/lib.d.ts",
      "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "signature": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8859c12c614ce56ba9a18e58384a198f-/// <reference no-default-lib=\"true\"/>\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array<T> { length: number; [n: number]: T; }\ninterface ReadonlyArray<T> {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "../node_modules/otherlib/index.d.ts",
      "version": "1fe659ed0634bb57b6dc25e9062f1162-export type Str = string;",
      "signature": "1fe659ed0634bb57b6dc25e9062f1162-export type Str = string;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../node_modules/somelib/index.d.ts",
      "version": "12e112ff6e2744bb42d8e0b511e44117-import type { Str } from 'otherLib';\nexport type Foo = { foo: Str; };",
      "signature": "12e112ff6e2744bb42d8e0b511e44117-import type { Str } from 'otherLib';\nexport type Foo = { foo: Str; };",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../node_modules/lib1/index.d.ts",
      "version": "b6305455d920a6729c435e6acf45eff6-import type { Foo } from 'someLib';\nexport type { Foo as Foo1 };",
      "signature": "b6305455d920a6729c435e6acf45eff6-import type { Foo } from 'someLib';\nexport type { Foo as Foo1 };",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../node_modules/lib2/index.d.ts",
      "version": "a5393e550a9c20a242a120bf6410db48-import type { Foo } from 'somelib';\nexport type { Foo as Foo2 };\nexport declare const foo2: Foo;",
      "signature": "a5393e550a9c20a242a120bf6410db48-import type { Foo } from 'somelib';\nexport type { Foo as Foo2 };\nexport declare const foo2: Foo;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./src/index.ts",
      "version": "42aef197ff5f079223e2c29fb2e77cc5-import type { Foo1 } from 'lib1';\nimport type { Foo2 } from 'lib2';\nexport const foo1: Foo1 = { foo: \"a\" };\nexport const foo2: Foo2 = { foo: \"b\" };",
      "signature": "42aef197ff5f079223e2c29fb2e77cc5-import type { Foo1 } from 'lib1';\nimport type { Foo2 } from 'lib2';\nexport const foo1: Foo1 = { foo: \"a\" };\nexport const foo2: Foo2 = { foo: \"b\" };",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "../node_modules/somelib/index.d.ts"
    ],
    [
      "../node_modules/otherlib/index.d.ts"
    ],
    [
      "../node_modules/lib1/index.d.ts",
      "../node_modules/lib2/index.d.ts"
    ]
  ],
  "referencedMap": {
    "../node_modules/lib1/index.d.ts": [
      "../node_modules/somelib/index.d.ts"
    ],
    "../node_modules/lib2/index.d.ts": [
      "../node_modules/somelib/index.d.ts"
    ],
    "../node_modules/somelib/index.d.ts": [
      "../node_modules/otherlib/index.d.ts"
    ],
    "./src/index.ts": [
      "../node_modules/lib1/index.d.ts",
      "../node_modules/lib2/index.d.ts"
    ]
  },
  "size": 1681
}
//// [/home/src/tslibs/TS/Lib/lib.d.ts] *Lib*
/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };

SemanticDiagnostics::
*refresh*    /home/src/tslibs/ts/lib/lib.d.ts
*refresh*    /home/node_modules/otherlib/index.d.ts
*refresh*    /home/node_modules/somelib/index.d.ts
*refresh*    /home/node_modules/lib1/index.d.ts
*refresh*    /home/node_modules/lib2/index.d.ts
*refresh*    /home/project/src/index.ts
Signatures::
