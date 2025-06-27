/**
 * 日付処理のユーティリティ関数
 */

/**
 * ISO形式の日付文字列をinput要素のdate型で使用できる形式(YYYY-MM-DD)に変換
 * タイムゾーンの影響を完全に除外して日付部分のみを抽出
 */
export const formatDateForInput = (dateString: string | undefined): string => {
  if (!dateString) return '';
  
  // 最初に文字列から直接YYYY-MM-DDパターンを抽出
  const datePattern = /^(\d{4})-(\d{2})-(\d{2})/;
  const match = dateString.match(datePattern);
  
  if (match) {
    // 文字列から直接抽出した場合はそのまま返す
    return `${match[1]}-${match[2]}-${match[3]}`;
  }
  
  // パターンにマッチしない場合は空文字列を返す
  return '';
};

/**
 * 日付を現在の環境のロケールで表示形式に変換
 */
export const formatDateForDisplay = (dateString: string | undefined): string => {
  if (!dateString) return '';
  try {
    const date = new Date(dateString);
    // 現在の環境のロケール設定を使用
    return date.toLocaleDateString();
  } catch {
    return '';
  }
};

/**
 * YYYY-MM-DD形式の日付をローカルタイムゾーンでISO形式に変換
 * @param dateStr YYYY-MM-DD形式の日付
 * @param isEndDate 終了日の場合は23:59:59を使用
 */
export const toISOStringWithJST = (dateStr: string, isEndDate: boolean = false): string => {
  if (!dateStr) return '';
  
  // 日付を検証
  const datePattern = /^(\d{4})-(\d{2})-(\d{2})$/;
  if (!datePattern.test(dateStr)) {
    return '';
  }
  
  // 現在のタイムゾーンオフセットを取得
  const now = new Date();
  const timezoneOffset = -now.getTimezoneOffset();
  const offsetHours = Math.floor(Math.abs(timezoneOffset) / 60);
  const offsetMinutes = Math.abs(timezoneOffset) % 60;
  const offsetSign = timezoneOffset >= 0 ? '+' : '-';
  const localOffset = `${offsetSign}${String(offsetHours).padStart(2, '0')}:${String(offsetMinutes).padStart(2, '0')}`;
  
  // 時刻部分を設定
  const time = isEndDate ? '23:59:59' : '00:00:00';
  
  // ローカルタイムゾーンでISO形式文字列を返す
  return `${dateStr}T${time}${localOffset}`;
};

/**
 * 安全な日付パース処理
 */
export const parseDate = (dateString: string | undefined): Date | null => {
  if (!dateString) return null;
  try {
    const date = new Date(dateString);
    // Invalid Dateチェック
    if (isNaN(date.getTime())) return null;
    return date;
  } catch {
    return null;
  }
};

/**
 * 日付の妥当性チェック
 */
export const isValidDateRange = (startDate: string, endDate: string): boolean => {
  const start = parseDate(startDate);
  const end = parseDate(endDate);
  
  if (!start || !end) return false;
  return start <= end;
};