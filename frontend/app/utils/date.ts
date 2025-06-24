/**
 * 日付処理のユーティリティ関数
 */

/**
 * ISO形式の日付文字列をinput要素のdate型で使用できる形式(YYYY-MM-DD)に変換
 * ローカルタイムゾーンを使用して変換
 */
export const formatDateForInput = (dateString: string | undefined): string => {
  if (!dateString) return '';
  try {
    const date = new Date(dateString);
    // ローカルタイムゾーンでの年月日を取得
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  } catch {
    return '';
  }
};

/**
 * 日付を日本語表示形式に変換
 */
export const formatDateForDisplay = (dateString: string | undefined): string => {
  if (!dateString) return '';
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP');
  } catch {
    return '';
  }
};

/**
 * YYYY-MM-DD形式の日付をJSTタイムゾーンを明示したISO形式に変換
 * @param dateStr YYYY-MM-DD形式の日付
 * @param isEndDate 終了日の場合は23:59:59を使用
 */
export const toISOStringWithJST = (dateStr: string, isEndDate: boolean = false): string => {
  if (!dateStr) return '';
  
  // JSTタイムゾーンオフセット
  const jstOffset = '+09:00';
  
  // 時刻部分を設定
  const time = isEndDate ? '23:59:59' : '00:00:00';
  
  // JSTタイムゾーンを明示したISO形式文字列を返す
  return `${dateStr}T${time}${jstOffset}`;
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