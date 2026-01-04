package test

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

// 方案洗牌算法性能测试

// ----------------------------------------------------------------
// 1. 基础结构与数据模拟
// ----------------------------------------------------------------

type DailyRecipePlan struct {
	MatchScore float64
	ID         int
}

// 模拟 Service 结构，持有 rng
type RecipeServiceMock struct {
	rng *rand.Rand
}

func NewService() *RecipeServiceMock {
	return &RecipeServiceMock{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ----------------------------------------------------------------
// 2. 待测试函数 (Implementation A: Fisher-Yates Variant)
// ----------------------------------------------------------------

// ShuffleInPlace (你的第一个实现：原地交换 + 预计算权重)
func (s *RecipeServiceMock) ShuffleInPlace(candidates []*DailyRecipePlan) []*DailyRecipePlan {
	n := len(candidates)
	if n == 0 {
		return nil
	}

	limit := min(n, 100)
	pool := candidates[:limit] // 注意：这里还是引用，原地修改会影响外部，Benchmark时需注意重置

	// 【优化点 1】预计算权重，Pow 只算 N 次
	weights := make([]float64, limit)
	for i, p := range pool {
		weights[i] = math.Pow(p.MatchScore, 3)
	}

	for i := 0; i < limit-1; i++ {
		remainingWeight := 0.0
		for j := i; j < limit; j++ {
			remainingWeight += weights[j]
		}

		r := s.rng.Float64() * remainingWeight

		currentSum := 0.0
		winnerIdx := -1

		for j := i; j < limit; j++ {
			currentSum += weights[j]
			if r <= currentSum {
				winnerIdx = j
				break
			}
		}

		if winnerIdx == -1 {
			winnerIdx = limit - 1
		}

		if winnerIdx != i {
			pool[i], pool[winnerIdx] = pool[winnerIdx], pool[i]
			weights[i], weights[winnerIdx] = weights[winnerIdx], weights[i]
		}
	}

	return pool
}

// ----------------------------------------------------------------
// 3. 待测试函数 (Implementation B: Construction/Selection)
// ----------------------------------------------------------------

// ShuffleConstruction (你的第二个实现：构建新切片 + 每次重新计算权重)
func (s *RecipeServiceMock) ShuffleConstruction(candidates []*DailyRecipePlan) []*DailyRecipePlan {
	if len(candidates) == 0 {
		return nil
	}

	topN := min(len(candidates), 100)
	pool := candidates[:topN]

	if len(pool) < 3 {
		return pool
	}

	// 复制一份 source，避免修改原数据
	source := make([]*DailyRecipePlan, len(pool))
	copy(source, pool)
	result := make([]*models_DailyRecipePlan, 0, len(pool)) // 注意：这里修正了你的代码，append需要相同类型

	// 【性能杀手预警】：在循环内
	for len(source) > 0 {
		totalWeight := 0.0
		// 【性能杀手 1】每次循环都重新遍历 source 计算总权重
		// 【性能杀手 2】每次循环都重复调用 math.Pow
		for _, p := range source {
			totalWeight += math.Pow(p.MatchScore, 3)
		}

		// 使用全局 rand (在并发下有锁竞争，单测中影响不大但也是开销)
		r := rand.Float64() * totalWeight
		curr := 0.0
		foundIdx := -1

		for i, p := range source {
			curr += math.Pow(p.MatchScore, 3)
			if r <= curr {
				foundIdx = i
				break
			}
		}

		if foundIdx == -1 {
			foundIdx = len(source) - 1
		}

		result = append(result, source[foundIdx])

		// 【性能杀手 3】切片删除操作，涉及内存移动
		source = append(source[:foundIdx], source[foundIdx+1:]...)
	}

	return result
}

// 辅助修正类型，为了代码能跑通
type models_DailyRecipePlan = DailyRecipePlan

// ----------------------------------------------------------------
// 4. Benchmark 测试代码
// ----------------------------------------------------------------

func makeCandidates(n int) []*DailyRecipePlan {
	res := make([]*DailyRecipePlan, n)
	for i := 0; i < n; i++ {
		res[i] = &DailyRecipePlan{
			ID:         i,
			MatchScore: rand.Float64() * 10.0, // 0-10 分
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 测试 Implementation A (InPlace)
func BenchmarkShuffle_InPlace_N100(b *testing.B) {
	s := NewService()
	// 准备足够的数据，每次 Reset 保证纯净（虽然 InPlace 只是打乱顺序，不影响下次运行的逻辑有效性）
	// 为了避免内存分配干扰计时，我们在外部创建好，但注意 InPlace 会修改顺序
	data := makeCandidates(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 为了公平对比（Implementation B 内部做了 copy），
		// 这里理论上 InPlace 不需要 copy，但为了防止数据被改得面目全非影响逻辑（虽然洗牌不影响），
		// 直接传原切片即可。InPlace 的优势就是 0 allocation (除了 weights 数组)。
		s.ShuffleInPlace(data)
	}
}

// 测试 Implementation B (Construction)
func BenchmarkShuffle_Construction_N100(b *testing.B) {
	s := NewService()
	data := makeCandidates(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ShuffleConstruction(data)
	}
}

// 增大压力测试 N=500 (虽然 limit 限制 100，但观察截断逻辑)
func BenchmarkShuffle_InPlace_N500(b *testing.B) {
	s := NewService()
	data := makeCandidates(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ShuffleInPlace(data)
	}
}

func BenchmarkShuffle_Construction_N500(b *testing.B) {
	s := NewService()
	data := makeCandidates(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ShuffleConstruction(data)
	}
}
