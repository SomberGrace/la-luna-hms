import { describe, expect, it } from 'vitest';
import { getBarItemImage, getRestaurantItemImage } from './productImages';

describe('product images', () => {
	it('uses a stored restaurant menu image when available', () => {
		expect(getRestaurantItemImage('Dish', 'Lunch', ' https://example.com/dish.jpg ')).toBe(
			'https://example.com/dish.jpg'
		);
	});

	it('provides a restaurant image fallback by name or category', () => {
		expect(getRestaurantItemImage('Grilled Lake Tilapia (Chambo)', 'Dinner')).toContain(
			'photo-1519708227418'
		);
		expect(getRestaurantItemImage('New item', 'Snack')).toContain('photo-1573080496219');
	});

	it.each([
		['Castle Lager (500ml)', 'photo-1608270119358'],
		['Coca-Cola 500ml', 'photo-1622483767028'],
		['Malawi Gin 250ml', 'photo-1514362545857']
	])('selects a matching bar product image for %s', (name, imageId) => {
		expect(getBarItemImage(name)).toContain(imageId);
	});
});
