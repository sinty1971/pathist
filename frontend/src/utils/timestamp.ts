import type { ModelsTimestamp } from '@/types/models';

/**
 * RFC3339Nano形式でサポートされる日時フォーマット（バックエンド対応）
 */
const TIMESTAMP_FORMATS_WITH_TZ = [
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}Z.*$/,     // RFC3339Nano
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z.*$/,            // RFC3339
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}.*$/ // RFC3339 with timezone
];

const TIMESTAMP_FORMATS_WITHOUT_TZ = [
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}.*$/,      // Without timezone
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*$/,             // Without timezone
  /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.*$/,             // Space separator
  /^\d{4}-\d{4}.*$/,                                       // YYYY-MMDD
  /^\d{4}-\d{2}-\d{2}.*$/,                                 // YYYY-MM-DD
  /^\d{8}.*$/,                                             // YYYYMMDD
  /^\d{4}\/\d{2}\/\d{2}.*$/,                               // YYYY/MM/DD
  /^\d{4}\.\d{2}\.\d{2}.*$/,                               // YYYY.MM.DD
  /^\d{4}\/\d{1,2}\/\d{1,2}.*$/,                          // YYYY/M/D
  /^\d{4}\.\d{1,2}\.\d{1,2}.*$/                           // YYYY.M.D
];

/**
 * 日時文字列を抽出するためのパターンマッチング
 */
function extractTimeString(input: string): { timeStr: string; remaining: string } | null {
  // タイムゾーン付きフォーマットを試行
  for (const regex of TIMESTAMP_FORMATS_WITH_TZ) {
    const match = input.match(regex);
    if (match) {
      const timeStr = match[0];
      const remaining = input.replace(timeStr, '').trim();
      return { timeStr, remaining };
    }
  }
  
  // タイムゾーンなしフォーマットを試行
  for (const regex of TIMESTAMP_FORMATS_WITHOUT_TZ) {
    const match = input.match(regex);
    if (match) {
      const timeStr = match[0];
      const remaining = input.replace(timeStr, '').trim();
      return { timeStr, remaining };
    }
  }
  
  return null;
}

/**
 * 文字列から日時をパースしてDateオブジェクトを返す（バックエンドのParseTime相当）
 * @param input 日時を含む文字列
 * @returns パースされたDateオブジェクトと残りの文字列
 */
export function parseTimeString(input: string): { date: Date; remaining: string } | null {
  if (!input) return null;
  
  const extracted = extractTimeString(input);
  if (!extracted) return null;
  
  const { timeStr, remaining } = extracted;
  
  // タイムゾーン付きフォーマットの試行
  const tzFormats = [
    (s: string) => new Date(s), // ISO形式（ブラウザが自動認識）
  ];
  
  for (const formatter of tzFormats) {
    try {
      const date = formatter(timeStr);
      if (!isNaN(date.getTime())) {
        return { date, remaining };
      }
    } catch (e) {
      continue;
    }
  }
  
  // タイムゾーンなしフォーマットの試行（JST として扱う）
  const noTzFormats = [
    // YYYY-MMDD 形式（例：2025-0618）
    (s: string) => {
      const match = s.match(/^(\d{4})-(\d{4})/);
      if (match) {
        const year = parseInt(match[1]);
        const monthDay = match[2];
        const month = parseInt(monthDay.substring(0, 2)) - 1; // 0-indexed
        const day = parseInt(monthDay.substring(2));
        return new Date(year, month, day);
      }
      return null;
    },
    // YYYY-MM-DD 形式
    (s: string) => {
      const match = s.match(/^(\d{4})-(\d{2})-(\d{2})/);
      if (match) {
        const year = parseInt(match[1]);
        const month = parseInt(match[2]) - 1;
        const day = parseInt(match[3]);
        return new Date(year, month, day);
      }
      return null;
    },
    // YYYYMMDD 形式
    (s: string) => {
      const match = s.match(/^(\d{8})/);
      if (match) {
        const dateStr = match[1];
        const year = parseInt(dateStr.substring(0, 4));
        const month = parseInt(dateStr.substring(4, 6)) - 1;
        const day = parseInt(dateStr.substring(6, 8));
        return new Date(year, month, day);
      }
      return null;
    },
    // その他の形式（時刻付き）
    (s: string) => {
      const isoStr = s.replace(' ', 'T');
      if (!isoStr.includes('T')) {
        return new Date(s);
      }
      // JSTとして扱う
      return new Date(isoStr + '+09:00');
    }
  ];
  
  for (const formatter of noTzFormats) {
    try {
      const date = formatter(timeStr);
      if (date && !isNaN(date.getTime())) {
        return { date, remaining };
      }
    } catch (e) {
      continue;
    }
  }
  
  return null;
}

/**
 * ModelsTimestampまたは文字列をISO文字列に変換
 */
export function timestampToString(timestamp?: ModelsTimestamp | string): string | undefined {
  if (!timestamp) {
    return undefined;
  }
  
  // 文字列の場合はそのまま返す（APIが直接文字列を返す場合）
  if (typeof timestamp === 'string') {
    return timestamp;
  }
  
  // オブジェクトの場合
  if (timestamp['time.Time']) {
    return timestamp['time.Time'];
  }
  
  return undefined;
}

/**
 * 文字列またはDateをModelsTimestampに変換（RFC3339Nano形式）
 */
export function stringToTimestamp(input?: string | Date): ModelsTimestamp | undefined {
  if (!input) {
    return undefined;
  }
  
  let date: Date;
  
  if (typeof input === 'string') {
    const parsed = parseTimeString(input);
    if (!parsed) {
      // フォールバック：直接Date constructor を使用
      date = new Date(input);
      if (isNaN(date.getTime())) {
        return undefined;
      }
    } else {
      date = parsed.date;
    }
  } else {
    date = input;
  }
  
  // RFC3339Nano形式で返す（バックエンドと同じフォーマット）
  const isoString = date.toISOString();
  // ナノ秒部分を追加（JavaScriptはミリ秒なので000を追加）
  const nanoString = isoString.replace(/\.\d{3}Z$/, '.000000000Z');
  
  return {
    'time.Time': nanoString
  };
}


/**
 * DateをRFC3339Nano形式の文字列に変換
 */
export function formatToRFC3339Nano(date: Date): string {
  const isoString = date.toISOString();
  return isoString.replace(/\.\d{3}Z$/, '.000000000Z');
}

/**
 * 現在の日時をJST（日本時間）で取得
 */
export function nowInJST(): Date {
  return new Date();
}

/**
 * 日時をJST（日本時間）の文字列として表示用にフォーマット
 */
export function formatDisplayDate(date: Date | string | ModelsTimestamp | undefined): string {
  if (!date) return '';
  
  let actualDate: Date;
  
  if (typeof date === 'string') {
    actualDate = new Date(date);
  } else if (date instanceof Date) {
    actualDate = date;
  } else {
    const dateStr = timestampToString(date);
    if (!dateStr) return '';
    actualDate = new Date(dateStr);
  }
  
  if (isNaN(actualDate.getTime())) return '';
  
  // 日本語ロケールで表示
  return actualDate.toLocaleString('ja-JP', {
    timeZone: 'Asia/Tokyo',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
}
