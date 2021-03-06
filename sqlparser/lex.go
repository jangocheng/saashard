// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The MIT License (MIT)

// Copyright (c) 2016 Jerry Bai

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sqlparser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/berkaroad/saashard/sqlparser/sqltypes"
)

// EOFCHAR EOF char.
const EOFCHAR = 0x100

// Tokenizer is the struct used to generate SQL
// tokens for the parser.
type Tokenizer struct {
	InStream      *strings.Reader
	AllowComments bool
	ForceEOF      bool
	lastChar      uint16
	Position      int
	errorToken    []byte
	LastError     string
	posVarIndex   int
	ParseTree     Statement
}

// NewStringTokenizer creates a new Tokenizer for the
// sql string.
func NewStringTokenizer(sql string) *Tokenizer {
	return &Tokenizer{InStream: strings.NewReader(sql)}
}

var keywords = map[string]int{
	"select": SELECT,
	"insert": INSERT,
	"update": UPDATE,
	"delete": DELETE,
	"from":   FROM,
	"where":  WHERE,
	"group":  GROUP,
	"having": HAVING,
	"order":  ORDER,
	"by":     BY,
	"limit":  LIMIT,
	"for":    FOR,

	"union":     UNION,
	"all":       ALL,
	"minus":     MINUS,
	"except":    EXCEPT,
	"intersect": INTERSECT,

	"join":          JOIN,
	"full":          FULL,
	"straight_join": STRAIGHT_JOIN,
	"left":          LEFT,
	"right":         RIGHT,
	"inner":         INNER,
	"outer":         OUTER,
	"cross":         CROSS,
	"natural":       NATURAL,
	"use":           USE,
	"force":         FORCE,
	"on":            ON,
	"into":          INTO,

	"distinct":  DISTINCT,
	"case":      CASE,
	"when":      WHEN,
	"then":      THEN,
	"else":      ELSE,
	"end":       END,
	"as":        AS,
	"and":       AND,
	"or":        OR,
	"not":       NOT,
	"exists":    EXISTS,
	"in":        IN,
	"is":        IS,
	"like":      LIKE,
	"between":   BETWEEN,
	"null":      NULL,
	"asc":       ASC,
	"desc":      DESC,
	"values":    VALUES,
	"duplicate": DUPLICATE,
	"key":       KEY,
	"default":   DEFAULT,
	"set":       SET,
	"lock":      LOCK,

	"create":   CREATE,
	"alter":    ALTER,
	"rename":   RENAME,
	"drop":     DROP,
	"table":    TABLE,
	"index":    INDEX,
	"view":     VIEW,
	"to":       TO,
	"ignore":   IGNORE,
	"if":       IF,
	"unique":   UNIQUE,
	"fulltext": FULLTEXT,
	"btree":    BTREE,
	"hash":     HASH,

	// Data Type
	"bit":        BIT,
	"tinyint":    TINYINT,
	"bool":       BOOL,
	"boolean":    BOOLEAN,
	"smallint":   SMALLINT,
	"mediumint":  MEDIUMINT,
	"int":        INT,
	"integer":    INTEGER,
	"bigint":     BIGINT,
	"real":       REAL,
	"double":     DOUBLE,
	"float":      FLOAT,
	"decimal":    DECIMAL,
	"date":       DATE,
	"time":       TIME,
	"timestamp":  TIMESTAMP,
	"datetime":   DATETIME,
	"year":       YEAR,
	"char":       CHAR,
	"nchar":      NCHAR,
	"varchar":    VARCHAR,
	"nvarchar":   NVARCHAR,
	"tinytext":   TINYTEXT,
	"text":       TEXT,
	"mediumtext": MEDIUMTEXT,
	"longtext":   LONGTEXT,
	"varbinary":  VARBINARY,
	"tinyblob":   TINYBLOB,
	"blob":       BLOB,
	"mediumblob": MEDIUMBLOB,
	"longblob":   LONGBLOB,

	"enum":           ENUM,
	"auto_increment": AUTO_INCREMENT,
	"engine":         ENGINE,
	"primary":        PRIMARY,
	"references":     REFERENCES,
	"comment":        COMMENT,

	"column_format": COLUMN_FORMAT,
	"fixed":         FIXED,
	"dynamic":       DYNAMIC,

	"disk":   DISK,
	"memory": MEMORY,

	"match":   MATCH,
	"partial": PARTIAL,
	"simple":  SIMPLE,

	"restrict": RESTRICT,
	"cascade":  CASCADE,

	"no":     NO,
	"action": ACTION,

	"unsigned": UNSIGNED,
	"zerofill": ZEROFILL,

	"constraint": CONSTRAINT,
	"foreign":    FOREIGN,

	"first": FIRST,
	"after": AFTER,

	"add":    ADD,
	"column": COLUMN,
	"change": CHANGE,
	"modify": MODIFY,

	"enable":  ENABLE,
	"disable": DISABLE,

	"using":    USING,
	"begin":    BEGIN,
	"rollback": ROLLBACK,
	"commit":   COMMIT,

	"names":        NAMES,
	"replace":      REPLACE,
	"start":        START,
	"transaction":  TRANSACTION,
	"isolation":    ISOLATION,
	"level":        LEVEL,
	"repeatable":   REPEATABLE,
	"read":         READ,
	"committed":    COMMITTED,
	"uncommitted":  UNCOMMITTED,
	"serializable": SERIALIZABLE,
	"collate":      COLLATE,
	"offset":       OFFSET,
	"charset":      CHARSET,
	"character":    CHARACTER,
	"collation":    COLLATION,
	"show":         SHOW,
	"describe":     DESCRIBE,
	"explain":      EXPLAIN,

	"variables":   VARIABLES,
	"status":      STATUS,
	"databases":   DATABASES,
	"database":    DATABASE,
	"tables":      TABLES,
	"columns":     COLUMNS,
	"fields":      FIELDS,
	"procedure":   PROCEDURE,
	"function":    FUNCTION,
	"engines":     ENGINES,
	"storage":     STORAGE,
	"plugins":     PLUGINS,
	"processlist": PROCESSLIST,
	"indexes":     INDEXES,
	"keys":        KEYS,
	"triggers":    TRIGGERS,
	"trigger":     TRIGGER,
	"slave":       SLAVE,

	"session": SESSION,
	"global":  GLOBAL,

	"profiles": PROFILES,

	// function
	"position": POSITION,

	"kill":       KILL,
	"query":      QUERY,
	"connection": CONNECTION,

	// charset
	"armscii8": ARMSCII8,
	"ascii":    ASCII,
	"big5":     BIG5,
	"binary":   BINARY,
	"cp1250":   CP1250,
	"cp1251":   CP1251,
	"cp1256":   CP1256,
	"cp1257":   CP1257,
	"cp850":    CP850,
	"cp852":    CP852,
	"cp866":    CP866,
	"cp932":    CP932,
	"dec8":     DEC8,
	"eucjpms":  EUCJPMS,
	"euckr":    EUCKR,
	"gb2312":   GB2312,
	"gbk":      GBK,
	"geostd8":  GEOSTD8,
	"greek":    GREEK,
	"hebrew":   HEBREW,
	"hp8":      HP8,
	"keybcs2":  KEYBCS2,
	"koi8r":    KOI8R,
	"koi8u":    KOI8U,
	"latin1":   LATIN1,
	"latin2":   LATIN2,
	"latin5":   LATIN5,
	"latin7":   LATIN7,
	"macce":    MACCE,
	"macroman": MACROMAN,
	"sjis":     SJIS,
	"swe7":     SWE7,
	"tis620":   TIS620,
	"ucs2":     UCS2,
	"ujis":     UJIS,
	"utf16":    UTF16,
	"utf16le":  UTF16LE,
	"utf32":    UTF32,
	"utf8":     UTF8,
	"utf8mb4":  UTF8MB4,

	// collate
	"armscii8_general_ci": ARMSCII8_GENERAL_CI,
	"armscii8_bin":        ARMSCII8_BIN,
	"ascii_general_ci":    ASCII_GENERAL_CI,
	"ascii_bin":           ASCII_BIN,
	"big5_chinese_ci":     BIG5_CHINESE_CI,
	"big5_bin":            BIG5_BIN,
	"cp1250_general_ci":   CP1250_GENERAL_CI,
	"cp1250_bin":          CP1250_BIN,
	"cp1251_chinese_ci":   CP1251_GENERAL_CI,
	"cp1251_chinese_cs":   CP1251_GENERAL_CS,
	"cp1251_bin":          CP1251_BIN,
	"cp1256_chinese_ci":   CP1256_GENERAL_CI,
	"cp1256_bin":          CP1256_BIN,
	"cp1257_chinese_ci":   CP1257_GENERAL_CI,
	"cp1257_bin":          CP1257_BIN,
	"cp850_chinese_ci":    CP850_GENERAL_CI,
	"cp850_bin":           CP850_BIN,
	"cp852_chinese_ci":    CP852_GENERAL_CI,
	"cp852_bin":           CP852_BIN,
	"cp866_chinese_ci":    CP866_GENERAL_CI,
	"cp866_bin":           CP866_BIN,
	"cp932_japanese_ci":   CP932_JAPANESE_CI,
	"cp932_bin":           CP932_BIN,
	"dec8_swedish_ci":     DEC8_SWEDISH_CI,
	"dec8_bin":            DEC8_BIN,
	"eucjpms_japanese_ci": EUCJPMS_JAPANESE_CI,
	"eucjpms_bin":         EUCJPMS_BIN,
	"euckr_korean_ci":     EUCKR_KOREAN_CI,
	"euckr_bin":           EUCKR_BIN,
	"gb2312_chinese_ci":   GB2312_CHINESE_CI,
	"gb2312_bin":          GB2312_BIN,
	"gbk_chinese_ci":      GBK_CHINESE_CI,
	"gbk_bin":             GBK_BIN,
	"geostd8_general_ci":  GEOSTD8_GENERAL_CI,
	"geostd8_bin":         GEOSTD8_BIN,
	"greek_general_ci":    GREEK_GENERAL_CI,
	"greek_bin":           GREEK_BIN,
	"hebrew_general_ci":   HEBREW_GENERAL_CI,
	"hebrew_bin":          HEBREW_BIN,
	"hp8_english_ci":      HP8_ENGLISH_CI,
	"hp8_bin":             HP8_BIN,
	"keybcs2_general_ci":  KEYBCS2_GENERAL_CI,
	"keybcs2_bin":         KEYBCS2_BIN,
	"koi8r_general_ci":    KOI8R_GENERAL_CI,
	"koi8r_bin":           KOI8R_BIN,
	"koi8u_general_ci":    KOI8U_GENERAL_CI,
	"koi8u_bin":           KOI8U_BIN,
	"latin1_general_ci":   LATIN1_GENERAL_CI,
	"latin1_general_cs":   LATIN1_GENERAL_CS,
	"latin1_bin":          LATIN1_BIN,
	"latin2_general_ci":   LATIN2_GENERAL_CI,
	"latin2_bin":          LATIN2_BIN,
	"latin5_turkish_ci":   LATIN5_TURKISH_CI,
	"latin5_bin":          LATIN5_BIN,
	"latin7_general_ci":   LATIN7_GENERAL_CI,
	"latin7_general_cs":   LATIN7_GENERAL_CS,
	"latin7_bin":          LATIN7_BIN,
	"macce_general_ci":    MACCE_GENERAL_CI,
	"macce_bin":           MACCE_BIN,
	"macroman_general_ci": MACROMAN_GENERAL_CI,
	"macroman_bin":        MACROMAN_BIN,
	"sjis_japanese_ci":    SJIS_JAPANESE_CI,
	"sjis_bin":            SJIS_BIN,
	"swe7_swedish_ci":     SWE7_SWEDISH_CI,
	"swe7_bin":            SWE7_BIN,
	"tis620_thai_ci":      TIS620_THAI_CI,
	"tis620_bin":          TIS620_BIN,
	"ucs2_general_ci":     UCS2_GENERAL_CI,
	"ucs2_unicode_ci":     UCS2_UNICODE_CI,
	"ucs2_bin":            UCS2_BIN,
	"ujis_japanese_ci":    UJIS_JAPANESE_CI,
	"ujis_bin":            UJIS_BIN,
	"utf16_general_ci":    UTF16_GENERAL_CI,
	"utf16_unicode_ci":    UTF16_UNICODE_CI,
	"utf16_bin":           UTF16_BIN,
	"utf16le_general_ci":  UTF16LE_GENERAL_CI,
	"utf16le_bin":         UTF16LE_BIN,
	"utf32_general_ci":    UTF32_GENERAL_CI,
	"utf32_unicode_ci":    UTF32_UNICODE_CI,
	"utf32_bin":           UTF32_BIN,
	"utf8_general_ci":     UTF8_GENERAL_CI,
	"utf8_unicode_ci":     UTF8_UNICODE_CI,
	"utf8_bin":            UTF8_BIN,
	"utf8mb4_general_ci":  UTF8MB4_GENERAL_CI,
	"utf8mb4_unicode_ci":  UTF8MB4_UNICODE_CI,
	"utf8mb4_bin":         UTF8MB4_BIN,
}

// Lex returns the next token form the Tokenizer.
// This function is used by go yacc.
func (tkn *Tokenizer) Lex(lval *yySymType) int {
	typ, val := tkn.Scan()
	for typ == COMMENTS {
		if tkn.AllowComments {
			break
		}
		typ, val = tkn.Scan()
	}
	switch typ {
	case ID, STRING, NUMBER, VALUE_ARG, COMMENTS:
		lval.bytes = val
	}
	tkn.errorToken = val
	return typ
}

// Error is called by go yacc if there's a parsing error.
func (tkn *Tokenizer) Error(err string) {
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	if tkn.errorToken != nil {
		fmt.Fprintf(buf, "%s at position %v near %s", err, tkn.Position, tkn.errorToken)
	} else {
		fmt.Fprintf(buf, "%s at position %v", err, tkn.Position)
	}
	tkn.LastError = buf.String()
}

// Scan scans the tokenizer for the next token and returns
// the token type and an optional value.
func (tkn *Tokenizer) Scan() (int, []byte) {
	if tkn.ForceEOF {
		return 0, nil
	}

	if tkn.lastChar == 0 {
		tkn.next()
	}
	tkn.skipBlank()
	switch ch := tkn.lastChar; {
	case isLetter(ch):
		return tkn.scanIdentifier()
	case isDigit(ch):
		return tkn.scanNumber(false)
	case ch == ':':
		return tkn.scanBindVar()
	default:
		tkn.next()
		switch ch {
		case EOFCHAR:
			return 0, nil
		case '=', ',', ';', '(', ')', '+', '*', '%', '&', '|', '^', '~':
			return int(ch), nil
		case '?':
			tkn.posVarIndex++
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, ":v%d", tkn.posVarIndex)
			return VALUE_ARG, []byte{'?'}
			// return VALUE_ARG, buf.Bytes()
		case '.':
			if isDigit(tkn.lastChar) {
				return tkn.scanNumber(true)
			}
			return int(ch), nil
		case '/':
			switch tkn.lastChar {
			case '/':
				tkn.next()
				return tkn.scanCommentType1("//")
			case '*':
				tkn.next()
				return tkn.scanCommentType2()
			default:
				return int(ch), nil
			}
		case '-':
			if tkn.lastChar == '-' {
				tkn.next()
				return tkn.scanCommentType1("--")
			}
			return int(ch), nil
		case '<':
			switch tkn.lastChar {
			case '>':
				tkn.next()
				return NE, nil
			case '=':
				tkn.next()
				switch tkn.lastChar {
				case '>':
					tkn.next()
					return NULL_SAFE_EQUAL, nil
				default:
					return LE, nil
				}
			default:
				return int(ch), nil
			}
		case '>':
			if tkn.lastChar == '=' {
				tkn.next()
				return GE, nil
			}
			return int(ch), nil
		case '!':
			if tkn.lastChar == '=' {
				tkn.next()
				return NE, nil
			}
			return LEX_ERROR, []byte("!")
		case '\'', '"':
			return tkn.scanString(ch, STRING)
		case '`':
			return tkn.scanString(ch, ID)
		default:
			return LEX_ERROR, []byte{byte(ch)}
		}
	}
}

// Next char.
func (tkn *Tokenizer) Next(buffer *bytes.Buffer) {
	if tkn.lastChar == EOFCHAR {
		// This should never happen.
		panic("unexpected EOF")
	}
	buffer.WriteByte(byte(tkn.lastChar))
	tkn.next()
}

func (tkn *Tokenizer) skipBlank() {
	ch := tkn.lastChar
	for isBlank(ch) {
		tkn.next()
		ch = tkn.lastChar
	}
}

func (tkn *Tokenizer) scanIdentifier() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteByte(byte(tkn.lastChar))
	for tkn.next(); isLetter(tkn.lastChar) || isDigit(tkn.lastChar); tkn.next() {
		buffer.WriteByte(byte(tkn.lastChar))
	}
	lowered := bytes.ToLower(buffer.Bytes())
	if keywordID, found := keywords[string(lowered)]; found {
		return keywordID, lowered
	}
	return ID, buffer.Bytes()
}

func (tkn *Tokenizer) scanBindVar() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteByte(byte(tkn.lastChar))
	for tkn.next(); isLetter(tkn.lastChar) || isDigit(tkn.lastChar) || tkn.lastChar == '.'; tkn.next() {
		buffer.WriteByte(byte(tkn.lastChar))
	}
	if buffer.Len() == 1 {
		return LEX_ERROR, buffer.Bytes()
	}
	return VALUE_ARG, buffer.Bytes()
}

func (tkn *Tokenizer) scanMantissa(base int, buffer *bytes.Buffer) {
	for digitVal(tkn.lastChar) < base {
		tkn.Next(buffer)
	}
}

func (tkn *Tokenizer) scanNumber(seenDecimalPoint bool) (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	if seenDecimalPoint {
		buffer.WriteByte('.')
		tkn.scanMantissa(10, buffer)
		goto exponent
	}

	if tkn.lastChar == '0' {
		// int or float
		tkn.Next(buffer)
		if tkn.lastChar == 'x' || tkn.lastChar == 'X' {
			// hexadecimal int
			tkn.Next(buffer)
			tkn.scanMantissa(16, buffer)
		} else {
			// octal int or float
			seenDecimalDigit := false
			tkn.scanMantissa(8, buffer)
			if tkn.lastChar == '8' || tkn.lastChar == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				tkn.scanMantissa(10, buffer)
			}
			if tkn.lastChar == '.' || tkn.lastChar == 'e' || tkn.lastChar == 'E' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				return LEX_ERROR, buffer.Bytes()
			}
		}
		goto exit
	}

	// decimal int or float
	tkn.scanMantissa(10, buffer)

fraction:
	if tkn.lastChar == '.' {
		tkn.Next(buffer)
		tkn.scanMantissa(10, buffer)
	}

exponent:
	if tkn.lastChar == 'e' || tkn.lastChar == 'E' {
		tkn.Next(buffer)
		if tkn.lastChar == '+' || tkn.lastChar == '-' {
			tkn.Next(buffer)
		}
		tkn.scanMantissa(10, buffer)
	}

exit:
	return NUMBER, buffer.Bytes()
}

func (tkn *Tokenizer) scanString(delim uint16, typ int) (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	for {
		ch := tkn.lastChar
		tkn.next()
		if ch == delim {
			if tkn.lastChar == delim {
				tkn.next()
			} else {
				break
			}
		} else if ch == '\\' {
			if tkn.lastChar == EOFCHAR {
				return LEX_ERROR, buffer.Bytes()
			}
			if decodedChar := sqltypes.SQLDecodeMap[byte(tkn.lastChar)]; decodedChar == sqltypes.DONTESCAPE {
				ch = tkn.lastChar
			} else {
				ch = uint16(decodedChar)
			}
			tkn.next()
		}
		if ch == EOFCHAR {
			return LEX_ERROR, buffer.Bytes()
		}
		buffer.WriteByte(byte(ch))
	}
	return typ, buffer.Bytes()
}

func (tkn *Tokenizer) scanCommentType1(prefix string) (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteString(prefix)
	for tkn.lastChar != EOFCHAR {
		if tkn.lastChar == '\n' {
			tkn.Next(buffer)
			break
		}
		tkn.Next(buffer)
	}
	return COMMENTS, buffer.Bytes()
}

func (tkn *Tokenizer) scanCommentType2() (int, []byte) {
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	buffer.WriteString("/*")
	for {
		if tkn.lastChar == '*' {
			tkn.Next(buffer)
			if tkn.lastChar == '/' {
				tkn.Next(buffer)
				break
			}
			continue
		}
		if tkn.lastChar == EOFCHAR {
			return LEX_ERROR, buffer.Bytes()
		}
		tkn.Next(buffer)
	}
	return COMMENTS, buffer.Bytes()
}

func (tkn *Tokenizer) next() {
	if ch, err := tkn.InStream.ReadByte(); err != nil {
		// Only EOF is possible.
		tkn.lastChar = EOFCHAR
	} else {
		tkn.lastChar = uint16(ch)
	}
	tkn.Position++
}

func isLetter(ch uint16) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '@'
}

func digitVal(ch uint16) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch) - '0'
	case 'a' <= ch && ch <= 'f':
		return int(ch) - 'a' + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch) - 'A' + 10
	}
	return 16 // larger than any legal digit val
}

func isDigit(ch uint16) bool {
	return '0' <= ch && ch <= '9'
}

func isBlank(ch uint16) bool {
	return ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t'
}
