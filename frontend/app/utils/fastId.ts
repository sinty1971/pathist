/**
 * 高速ID生成 - Go版の移植（最適化済み）
 * 元のGo実装から約4-8倍遅くなる予想
 */

// BLAKE2bハッシュ用（軽量実装を使用）
import { blake2b } from 'blake2b-wasm';

// 32進数変換用文字テーブル（Go版と同じ）
const RADIX_TABLE = '123456789ABCDEFGHJKLMNPRSTUVWXYZ';
const BASE = 32n;

/**
 * 高速ID生成クラス
 */
export class FastID {
  private hash128: bigint;

  constructor(hash128: bigint) {
    this.hash128 = hash128;
  }

  /**
   * バイト配列からIDを生成（最適化版）
   */
  static fromBytes(data: Uint8Array): FastID {
    // BLAKE2b-256を使用（Go版と同等）
    const hash = blake2b(32).update(data).digest();
    
    // 下位128ビットを使用（Go版と同じ）
    const lower128Bytes = hash.slice(16);
    
    // BigIntに変換（リトルエンディアン）
    let hash128 = 0n;
    for (let i = lower128Bytes.length - 1; i >= 0; i--) {
      hash128 = hash128 * 256n + BigInt(lower128Bytes[i]);
    }

    return new FastID(hash128);
  }

  /**
   * 文字列からIDを生成（最適化版）
   */
  static fromString(str: string): FastID {
    // TextEncoderを使用してUTF-8バイト配列に変換
    const encoder = new TextEncoder();
    return FastID.fromBytes(encoder.encode(str));
  }

  /**
   * 工事ID生成（Goの GenerateKojiID 相当）
   */
  static generateKojiID(
    startDate: Date,
    companyName: string,
    locationName: string
  ): string {
    // Goと同じフォーマットで文字列を構築
    const dateStr = startDate.toISOString().slice(0, 10).replace(/-/g, '');
    const combined = `${dateStr}${companyName}${locationName}`;
    
    return FastID.fromString(combined).len5();
  }

  /**
   * 5文字のIDを返す（32^5 = 約3,355万通り）- 最適化版
   */
  len5(): string {
    return this.toBase32(5);
  }

  /**
   * 7文字のIDを返す（32^7 = 約34億通り）
   */
  len7(): string {
    return this.toBase32(7);
  }

  /**
   * 25文字のIDを返す（32^25 = 約10^37通り）
   */
  full25(): string {
    return this.toBase32(25);
  }

  /**
   * 32進数変換（最適化版）
   */
  private toBase32(length: number): string {
    const bytes = new Array(length);
    let value = this.hash128;

    // Go版と同じアルゴリズム（最下位桁から計算）
    for (let i = length - 1; i >= 0; i--) {
      const mod = value % BASE;
      value = value / BASE;
      bytes[i] = RADIX_TABLE[Number(mod)];
    }

    return bytes.join('');
  }

  /**
   * デフォルトで5文字のIDを返す
   */
  toString(): string {
    return this.len5();
  }
}

/**
 * 工事フォルダ名生成（Go版の移植）
 */
export function generateKojiFolderName(
  startDate: Date,
  companyName: string,
  locationName: string
): string {
  if (!startDate || isNaN(startDate.getTime())) {
    throw new Error('開始日が無効です');
  }

  const year = startDate.getFullYear();
  const month = (startDate.getMonth() + 1).toString().padStart(2, '0');
  const day = startDate.getDate().toString().padStart(2, '0');

  return `${year}-${month}${day} ${companyName} ${locationName}`;
}

/**
 * 使用例とベンチマーク用関数
 */
export function benchmarkIdGeneration(iterations: number = 10000): number {
  const startTime = performance.now();
  
  for (let i = 0; i < iterations; i++) {
    FastID.generateKojiID(
      new Date(2025, 5, 18),
      '豊田築炉',
      '名和工場'
    );
  }
  
  const endTime = performance.now();
  const totalMs = endTime - startTime;
  const nsPerOp = (totalMs * 1_000_000) / iterations;
  
  console.log(`TypeScript ID生成ベンチマーク:`);
  console.log(`${iterations}回実行: ${totalMs.toFixed(2)}ms`);
  console.log(`1回あたり: ${nsPerOp.toFixed(0)}ns/op`);
  console.log(`Goと比較: 約${Math.round(nsPerOp / 100)}倍遅い（推定）`);
  
  return nsPerOp;
}