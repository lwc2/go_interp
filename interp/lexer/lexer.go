package lexer

import (
	"go_interp/model/token"
	"go_interp/util"
)

type Lexer struct {
	input string
	// 所输入字符串中的当前位置(指向当前字符串)
	position int
	// 所输入字符串中的当前读取位置(指向当前字符只有的一个字符)
	readPosition int
	ch           byte
}

func Load(input string) *Lexer {
	I := &Lexer{
		input: input,
	}
	I.readChar()
	return I
}

func (I *Lexer) readChar() {
	if I.readPosition >= len(I.input) {
		// 0是NUL字符的ASCII编码, 用来表示"尚未读取任何内容" 或 EOF
		I.ch = 0
	} else {
		I.ch = I.input[I.readPosition]
	}
	I.position = I.readPosition
	I.readPosition += 1
}

func (I *Lexer) peekChar() byte {
	if I.readPosition >= len(I.input) {
		return 0
	} else {
		return I.input[I.readPosition]
	}
}

func (I *Lexer) readIdentifier() (str string) {
	position := I.position
	for util.IsLetter(I.ch) {
		I.readChar()
	}
	return I.input[position:I.position]
}

func (I *Lexer) NextToken() (tok token.Token) {
	I.skipWhiteSpace()

	switch I.ch {
	case '=':
		if I.peekChar() == '=' {
			ch := I.ch
			I.readChar()
			lit := string(ch) + string(I.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: lit,
			}
		} else {
			tok = token.NewToken(token.ASSIGN, I.ch)
		}
	case '+':
		tok = token.NewToken(token.PLUS, I.ch)
	case '-':
		tok = token.NewToken(token.MINUS, I.ch)
	case '!':
		if I.peekChar() == '=' {
			ch := I.ch
			I.readChar()
			lit := string(ch) + string(I.ch)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: lit,
			}
		} else {
			tok = token.NewToken(token.BANG, I.ch)
		}
	case '<':
		tok = token.NewToken(token.LT, I.ch)
	case '>':
		tok = token.NewToken(token.GT, I.ch)
	case '*':
		tok = token.NewToken(token.ASTERISK, I.ch)
	case '/':
		tok = token.NewToken(token.SLASH, I.ch)
	case '(':
		tok = token.NewToken(token.LPAREN, I.ch)
	case ')':
		tok = token.NewToken(token.RPAREN, I.ch)
	case '{':
		tok = token.NewToken(token.LBRACE, I.ch)
	case '}':
		tok = token.NewToken(token.RBRACE, I.ch)
	case ';':
		tok = token.NewToken(token.SEMICOLON, I.ch)
	case ',':
		tok = token.NewToken(token.COMMA, I.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if util.IsLetter(I.ch) {
			// 处理关键字和变量
			tok.Literal = I.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return
		} else if util.IsDigit(I.ch) {
			// 处理数字
			tok.Type = token.INT
			tok.Literal = I.readNumber()
			return
		} else {
			// 处理错误
			tok = token.NewToken(token.ILLEGAL, I.ch)
		}
	}

	I.readChar()
	return
}

func (I *Lexer) skipWhiteSpace() {
	for I.ch == ' ' || I.ch == '\t' || I.ch == '\n' || I.ch == '\r' {
		I.readChar()
	}
}

func (I *Lexer) readNumber() string {
	position := I.position
	for util.IsDigit(I.ch) {
		I.readChar()
	}
	return I.input[position:I.position]
}
