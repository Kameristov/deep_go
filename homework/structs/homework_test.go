package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем 32-й бит (флаг дома) и устанавливаем золото в первые 31 бит
		person.goldHome = (person.goldHome & 0x80000000) | uint32(gold&0x7FFFFFFF)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем биты 11-16 и устанавливаем ману в первые 10 бит
		person.manaRespectGunFamily = (person.manaRespectGunFamily & 0xFC00) | uint16(mana&0x3FF)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем биты 11-16 и устанавливаем здоровье в первые 10 бит
		person.healthStrengthType = (person.healthStrengthType & 0xFC00) | uint16(health&0x3FF)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем биты 1-10 и 15-16, устанавливаем уважение в биты 11-14
		person.manaRespectGunFamily = (person.manaRespectGunFamily & 0xC3FF) | (uint16(respect&0x0F) << 10)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем биты 1-10 и 15-16, устанавливаем силу в биты 11-14
		person.healthStrengthType = (person.healthStrengthType & 0xC3FF) | (uint16(strength&0x0F) << 10)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем биты 5-8 (уровень) и устанавливаем опыт в первые 4 бита
		person.experienceLevel = (person.experienceLevel & 0xF0) | uint8(experience&0x0F)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем первые 4 бита (опыт) и устанавливаем уровень в последние 4 бита
		person.experienceLevel = (person.experienceLevel & 0x0F) | (uint8(level&0x0F) << 4)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		// Устанавливаем 1 в 32-м бите (0x80000000 = 1000 0000 0000 0000 0000 0000 0000 0000)
		person.goldHome |= 0x80000000
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		// Устанавливаем 1 в 15-м бите (0x4000 = 0100 0000 0000 0000)
		person.manaRespectGunFamily |= 0x4000
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		// Устанавливаем 1 в 16-м бите (0x8000 = 1000 0000 0000 0000)
		person.manaRespectGunFamily |= 0x8000
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		// Сохраняем остальные биты, устанавливаем тип в последние 2 бита
		person.healthStrengthType = (person.healthStrengthType & 0xFFFC) | uint16(personType&0x03)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

//  Имя пользователя [0…42] символов латиницы - {нужно 42 байта} +
//  Координата по оси X [-2_000_000_000…2_000_000_000] значений - {нужно 4 байта}
//  Координата по оси Y [-2_000_000_000…2_000_000_000] значений - {нужно 4 байта}
//  Координата по оси Z [-2_000_000_000…2_000_000_000] значений - {нужно 4 байта}
//  Золото [0…2_000_000_000] значений - {нужно 31 бит}
//  Магическая сила (мана) [0…1000] значений - {нужно 10 бит}
//  Здоровье [0…1000] значений - {нужно 10 бит}
//  Уважение [0…10] значений - {нужно 4 бит}
//  Сила [0…10] значений - {нужно 4 бит}
//  Опыт [0…10] значений - {нужно 4 бит}
//  Уровень [0…10] значений - {нужно 4 бит}
//  Есть ли у игрока дом [true/false] значения - {нужно 1 бит}
//  Есть ли у игрока оружие [true/false] значения - {нужно 1 бит}
//  Есть ли у игрока семья [true/false] значения - {нужно 1 бит}
//  Тип игрока [строитель/кузнец/воин] значения - {нужно 2 бит}

type GamePerson struct {
	x                    int32
	y                    int32
	z                    int32
	goldHome             uint32 // Золото31 + дом1
	manaRespectGunFamily uint16 // мана10 + Уважение4 + оружие1 + семья1
	healthStrengthType   uint16 // Здоровье10 + Сила4 + Тип2
	experienceLevel      uint8  // Опыт4 + Уровень4
	name                 [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.goldHome & 0x7FFFFFFF)
}

func (p *GamePerson) Mana() int {
	return int(p.manaRespectGunFamily & 0x3FF)
}

func (p *GamePerson) Health() int {
	return int(p.healthStrengthType & 0x3FF)
}

func (p *GamePerson) Respect() int {
	return int((p.manaRespectGunFamily >> 10) & 0x0F)
}

func (p *GamePerson) Strength() int {
	return int((p.healthStrengthType >> 10) & 0x0F)
}

func (p *GamePerson) Experience() int {
	return int(p.experienceLevel & 0x0F)
}

func (p *GamePerson) Level() int {
	return int((p.experienceLevel >> 4) & 0x0F)
}

func (p *GamePerson) HasHouse() bool {
	return (p.goldHome & 0x80000000) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.manaRespectGunFamily & 0x4000) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return (p.manaRespectGunFamily & 0x8000) != 0
}

func (p *GamePerson) Type() int {
	return int(p.healthStrengthType & 0x03)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
